// Package scheduler runs the daily GDHR sync window (default 18:00 -> 05:00 in
// SYNC_TIMEZONE). Within a window it launches the per-entity sync workers in
// parallel; at window end (or shutdown) their context is cancelled. The durable
// Redis cursor means a restart mid-window resumes where it left off.
package scheduler

import (
	"context"
	"log"
	"strconv"
	"strings"
	"sync"
	"time"

	"backend/config"
	syncsvc "backend/services/sync"
)

// Scheduler owns the timing loop and the active run's lifecycle.
type Scheduler struct {
	cfg *config.Config
	svc *syncsvc.Service
	loc *time.Location

	startMin int // minutes from midnight
	endMin   int

	mu        sync.Mutex
	running   bool
	activeDay string
	cancel    context.CancelFunc
}

// New builds a Scheduler. An unparseable timezone falls back to local time.
func New(cfg *config.Config, svc *syncsvc.Service) (*Scheduler, error) {
	loc, err := time.LoadLocation(cfg.SyncTimezone)
	if err != nil {
		log.Printf("[sched] timezone %q not found, using local: %v", cfg.SyncTimezone, err)
		loc = time.Local
	}
	return &Scheduler{
		cfg:      cfg,
		svc:      svc,
		loc:      loc,
		startMin: parseHM(cfg.SyncWindowStart, 18*60),
		endMin:   parseHM(cfg.SyncWindowEnd, 5*60),
	}, nil
}

// Start launches the background timing loop.
func (s *Scheduler) Start(ctx context.Context) {
	go s.loop(ctx)
}

func (s *Scheduler) loop(ctx context.Context) {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	s.tick(ctx)
	for {
		select {
		case <-ctx.Done():
			s.mu.Lock()
			if s.cancel != nil {
				s.cancel()
			}
			s.mu.Unlock()
			return
		case <-ticker.C:
			s.tick(ctx)
		}
	}
}

// tick starts a run when the clock first enters the daily window.
func (s *Scheduler) tick(parent context.Context) {
	now := time.Now().In(s.loc)
	if !s.inWindow(now) {
		return
	}
	day := s.runDate(now)

	s.mu.Lock()
	skip := s.running || s.activeDay == day
	s.mu.Unlock()
	if skip {
		return
	}

	s.startRun(parent, day, s.windowEnd(day), true)
}

// TriggerNow starts an on-demand run for today, ignoring the window. Returns
// false if a run is already active. Used by the manual-trigger endpoint.
func (s *Scheduler) TriggerNow(parent context.Context) (bool, string) {
	day := time.Now().In(s.loc).Format("2006-01-02")
	if s.startRun(parent, day, time.Time{}, false) {
		return true, day
	}
	return false, day
}

// startRun launches runDay under a fresh context, guarded so only one run is
// ever active. With hasDeadline the run is auto-cancelled at the window end.
func (s *Scheduler) startRun(parent context.Context, day string, deadline time.Time, hasDeadline bool) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.running {
		return false
	}

	var runCtx context.Context
	var cancel context.CancelFunc
	if hasDeadline {
		runCtx, cancel = context.WithDeadline(parent, deadline)
	} else {
		runCtx, cancel = context.WithCancel(parent)
	}

	s.running = true
	s.activeDay = day
	s.cancel = cancel

	go func() {
		s.runDay(runCtx, day)
		cancel()
		s.mu.Lock()
		s.running = false
		s.mu.Unlock()
	}()
	return true
}

// runDay launches the four sync workers in parallel for the given run date.
func (s *Scheduler) runDay(ctx context.Context, day string) {
	log.Printf("[sched] sync run starting for %s", day)

	var wg sync.WaitGroup
	launch := func(f func()) {
		wg.Add(1)
		go func() {
			defer wg.Done()
			f()
		}()
	}

	launch(func() { s.svc.SyncRanks(ctx, day) })
	launch(func() { s.svc.SyncPositions(ctx, day) })
	launch(func() { s.svc.SyncInstitutes(ctx, day, s.cfg.SyncInstitutesInterval) })
	launch(func() { s.svc.SyncStaffs(ctx, day, s.cfg.SyncStaffsInterval) })

	wg.Wait()
	log.Printf("[sched] sync run finished/stopped for %s", day)
}

// Today returns the current date string in the scheduler's timezone.
func (s *Scheduler) Today() string {
	return time.Now().In(s.loc).Format("2006-01-02")
}

// inWindow reports whether now's time-of-day is inside the window (wrap-aware).
func (s *Scheduler) inWindow(now time.Time) bool {
	cur := now.Hour()*60 + now.Minute()
	if s.startMin <= s.endMin {
		return cur >= s.startMin && cur < s.endMin
	}
	return cur >= s.startMin || cur < s.endMin // wraps past midnight
}

// runDate returns the date the window's start belongs to. For a wrap-around
// window, the after-midnight tail belongs to the previous day.
func (s *Scheduler) runDate(now time.Time) string {
	cur := now.Hour()*60 + now.Minute()
	if s.startMin > s.endMin && cur < s.startMin {
		return now.AddDate(0, 0, -1).Format("2006-01-02")
	}
	return now.Format("2006-01-02")
}

// windowEnd returns the absolute end time for the window that started on day.
func (s *Scheduler) windowEnd(day string) time.Time {
	d, err := time.ParseInLocation("2006-01-02", day, s.loc)
	if err != nil {
		return time.Now().In(s.loc).Add(11 * time.Hour)
	}
	end := time.Date(d.Year(), d.Month(), d.Day(), s.endMin/60, s.endMin%60, 0, 0, s.loc)
	if s.endMin <= s.startMin {
		end = end.AddDate(0, 0, 1) // window crosses midnight
	}
	return end
}

// parseHM parses "HH:MM" into minutes-from-midnight, returning fallback on error.
func parseHM(v string, fallback int) int {
	parts := strings.SplitN(strings.TrimSpace(v), ":", 2)
	if len(parts) != 2 {
		return fallback
	}
	h, err1 := strconv.Atoi(strings.TrimSpace(parts[0]))
	m, err2 := strconv.Atoi(strings.TrimSpace(parts[1]))
	if err1 != nil || err2 != nil || h < 0 || h > 23 || m < 0 || m > 59 {
		return fallback
	}
	return h*60 + m
}

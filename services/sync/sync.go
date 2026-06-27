// Package sync imports GDHR data into MySQL. Paged entities (institutes, staffs)
// walk page-by-page, persisting a durable cursor in Redis so a restart resumes
// instead of starting over. Every write is an UPSERT, so re-running refreshes
// existing rows rather than duplicating them.
package sync

import (
	"context"
	"log"
	"time"

	"backend/models"
	"backend/services/gdhr"
	"backend/services/hierarchy"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// Service performs the imports.
type Service struct {
	client *gdhr.Client
	db     *gorm.DB
}

// New builds a sync Service.
func New(client *gdhr.Client, db *gorm.DB) *Service {
	return &Service{client: client, db: db}
}

const (
	maxRetries     = 3
	retryBackoff   = 5 * time.Second
	insertBatchMax = 200
)

// upsertAll writes rows with ON CONFLICT UPDATE ALL (insert-or-update by PK).
func (s *Service) upsertAll(rows any) error {
	return s.db.Clauses(clause.OnConflict{UpdateAll: true}).
		CreateInBatches(rows, insertBatchMax).Error
}

// ---- institutes (paged, interval-driven) ----

func (s *Service) SyncInstitutes(ctx context.Context, date string, interval time.Duration) {
	const entity = "institutes"

	if done, err := IsDone(ctx, entity, date); err != nil {
		log.Printf("[sync] %s cursor read failed: %v", entity, err)
		return
	} else if done {
		log.Printf("[sync] %s already complete for %s", entity, date)
		return
	}

	for {
		if ctx.Err() != nil {
			log.Printf("[sync] %s stopped (window/shutdown) for %s", entity, date)
			return
		}

		page, err := NextPage(ctx, entity, date)
		if err != nil {
			log.Printf("[sync] %s cursor read failed: %v", entity, err)
			return
		}

		resp, err := s.fetchInstitutes(ctx, page)
		if err != nil {
			log.Printf("[sync] %s page %d fetch failed (cursor not advanced): %v", entity, page, err)
			return
		}

		if len(resp.Items) == 0 {
			if err := MarkDone(ctx, entity, date); err != nil {
				log.Printf("[sync] %s mark-done failed: %v", entity, err)
			}
			log.Printf("[sync] %s complete at page %d for %s", entity, page, date)
			return
		}

		rows := make([]models.Institute, 0, len(resp.Items))
		for _, d := range resp.Items {
			rows = append(rows, mapInstitute(d))
		}
		if err := s.upsertAll(&rows); err != nil {
			log.Printf("[sync] %s page %d upsert failed (cursor not advanced): %v", entity, page, err)
			return
		}

		if err := SetNextPage(ctx, entity, date, page+1); err != nil {
			log.Printf("[sync] %s cursor advance failed: %v", entity, err)
			return
		}
		log.Printf("[sync] %s page %d upserted %d rows", entity, page, len(rows))

		if !wait(ctx, interval) {
			return
		}
	}
}

// ---- staffs (paged, interval-driven) ----

func (s *Service) SyncStaffs(ctx context.Context, date string, interval time.Duration) {
	const entity = "staffs"

	if done, err := IsDone(ctx, entity, date); err != nil {
		log.Printf("[sync] %s cursor read failed: %v", entity, err)
		return
	} else if done {
		log.Printf("[sync] %s already complete for %s", entity, date)
		return
	}

	// Build the institute tree once so each staff can be tagged with its full
	// institute chain (ids + names) without per-row queries.
	ix, err := hierarchy.BuildIndex(s.db)
	if err != nil {
		log.Printf("[sync] %s: institute index build failed, chains will be empty: %v", entity, err)
	}

	for {
		if ctx.Err() != nil {
			log.Printf("[sync] %s stopped (window/shutdown) for %s", entity, date)
			return
		}

		page, err := NextPage(ctx, entity, date)
		if err != nil {
			log.Printf("[sync] %s cursor read failed: %v", entity, err)
			return
		}

		resp, err := s.fetchStaffs(ctx, page)
		if err != nil {
			log.Printf("[sync] %s page %d fetch failed (cursor not advanced): %v", entity, page, err)
			return
		}

		if len(resp.Data) == 0 {
			if err := MarkDone(ctx, entity, date); err != nil {
				log.Printf("[sync] %s mark-done failed: %v", entity, err)
			}
			log.Printf("[sync] %s complete at page %d for %s", entity, page, date)
			return
		}

		rows := make([]models.Staff, 0, len(resp.Data))
		for _, d := range resp.Data {
			st := mapStaff(d)
			if ix != nil {
				st.InstituteIDs, st.InstituteHierarchy = ix.Chain(st.InstituteID)
			}
			rows = append(rows, st)
		}
		if err := s.upsertAll(&rows); err != nil {
			log.Printf("[sync] %s page %d upsert failed (cursor not advanced): %v", entity, page, err)
			return
		}

		if err := SetNextPage(ctx, entity, date, page+1); err != nil {
			log.Printf("[sync] %s cursor advance failed: %v", entity, err)
			return
		}
		log.Printf("[sync] %s page %d upserted %d rows", entity, page, len(rows))

		if !wait(ctx, interval) {
			return
		}
	}
}

// ---- ranks (single call, once per day) ----

func (s *Service) SyncRanks(ctx context.Context, date string) {
	const entity = "ranks"

	if done, err := IsDone(ctx, entity, date); err == nil && done {
		log.Printf("[sync] %s already complete for %s", entity, date)
		return
	}

	resp, err := s.fetchRanks(ctx)
	if err != nil {
		log.Printf("[sync] %s fetch failed: %v", entity, err)
		return
	}
	if len(resp.Data.Items) == 0 {
		log.Printf("[sync] %s returned no rows", entity)
		_ = MarkDone(ctx, entity, date)
		return
	}

	rows := make([]models.Rank, 0, len(resp.Data.Items))
	for _, d := range resp.Data.Items {
		rows = append(rows, mapRank(d))
	}
	if err := s.upsertAll(&rows); err != nil {
		log.Printf("[sync] %s upsert failed: %v", entity, err)
		return
	}
	_ = MarkDone(ctx, entity, date)
	log.Printf("[sync] %s upserted %d rows for %s", entity, len(rows), date)
}

// ---- positions (single call, once per day) ----

func (s *Service) SyncPositions(ctx context.Context, date string) {
	const entity = "positions"

	if done, err := IsDone(ctx, entity, date); err == nil && done {
		log.Printf("[sync] %s already complete for %s", entity, date)
		return
	}

	resp, err := s.fetchPositions(ctx)
	if err != nil {
		log.Printf("[sync] %s fetch failed: %v", entity, err)
		return
	}
	if len(resp.Data.Items) == 0 {
		log.Printf("[sync] %s returned no rows", entity)
		_ = MarkDone(ctx, entity, date)
		return
	}

	rows := make([]models.Position, 0, len(resp.Data.Items))
	for _, d := range resp.Data.Items {
		rows = append(rows, mapPosition(d))
	}
	if err := s.upsertAll(&rows); err != nil {
		log.Printf("[sync] %s upsert failed: %v", entity, err)
		return
	}
	_ = MarkDone(ctx, entity, date)
	log.Printf("[sync] %s upserted %d rows for %s", entity, len(rows), date)
}

// ---- fetch-with-retry wrappers ----

func (s *Service) fetchInstitutes(ctx context.Context, page int) (*gdhr.InstitutesResponse, error) {
	var resp *gdhr.InstitutesResponse
	err := retry(ctx, func() error {
		r, e := s.client.FetchInstitutes(ctx, page)
		resp = r
		return e
	})
	return resp, err
}

func (s *Service) fetchStaffs(ctx context.Context, page int) (*gdhr.StaffsResponse, error) {
	var resp *gdhr.StaffsResponse
	err := retry(ctx, func() error {
		r, e := s.client.FetchStaffs(ctx, page)
		resp = r
		return e
	})
	return resp, err
}

func (s *Service) fetchRanks(ctx context.Context) (*gdhr.RanksResponse, error) {
	var resp *gdhr.RanksResponse
	err := retry(ctx, func() error {
		r, e := s.client.FetchRanks(ctx)
		resp = r
		return e
	})
	return resp, err
}

func (s *Service) fetchPositions(ctx context.Context) (*gdhr.PositionsResponse, error) {
	var resp *gdhr.PositionsResponse
	err := retry(ctx, func() error {
		r, e := s.client.FetchPositions(ctx)
		resp = r
		return e
	})
	return resp, err
}

// retry runs fn up to maxRetries times with a fixed backoff, aborting on ctx.
func retry(ctx context.Context, fn func() error) error {
	var err error
	for attempt := 1; attempt <= maxRetries; attempt++ {
		if err = fn(); err == nil {
			return nil
		}
		if ctx.Err() != nil {
			return ctx.Err()
		}
		if attempt < maxRetries {
			if !wait(ctx, retryBackoff) {
				return ctx.Err()
			}
		}
	}
	return err
}

// wait sleeps for d, returning false if ctx is cancelled first.
func wait(ctx context.Context, d time.Duration) bool {
	t := time.NewTimer(d)
	defer t.Stop()
	select {
	case <-ctx.Done():
		return false
	case <-t.C:
		return true
	}
}

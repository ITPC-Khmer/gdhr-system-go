// Command backfill runs a one-off GDHR import in the foreground (no scheduler
// window, no 20s pacing), looping page 1 -> end for the chosen entity. It reuses
// the same client + upsert + Redis cursor as the scheduled sync, so it resumes
// on restart and is safe to re-run (upsert).
//
// Usage:
//
//	go run ./cmd/backfill -entity=institutes -interval=500ms -fresh
//
// Credentials come from the environment / .env (GDHR_USER, GDHR_PASS, GDHR_KEY).
package main

import (
	"context"
	"flag"
	"log"
	"time"

	"backend/cache"
	"backend/config"
	"backend/database"
	"backend/models"
	"backend/services/gdhr"
	"backend/services/hierarchy"
	syncsvc "backend/services/sync"
)

func main() {
	entity := flag.String("entity", "institutes", "institutes | staffs | ranks | positions | all | staff-hierarchy")
	interval := flag.Duration("interval", 500*time.Millisecond, "delay between pages (be polite to the API)")
	fresh := flag.Bool("fresh", false, "clear the Redis cursor and start from page 1")
	flag.Parse()

	cfg := config.Load()
	database.Connect(cfg)

	// staff-hierarchy is DB-only (no API, no Redis): recompute institute chains
	// for all existing staff, grouped by distinct institute_id for efficiency.
	if *entity == "staff-hierarchy" {
		backfillStaffHierarchy()
		return
	}

	if cfg.GDHRUser == "" || cfg.GDHRKey == "" {
		log.Fatal("GDHR credentials missing — set GDHR_USER / GDHR_PASS / GDHR_KEY in env or .env")
	}

	cache.Connect(cfg)
	if cache.RDB == nil {
		log.Fatal("redis unavailable — required for the page cursor")
	}

	client := gdhr.New(cfg.GDHRBaseURL, cfg.GDHRUser, cfg.GDHRPass, cfg.GDHRKey, cfg.SyncHTTPTimeout)
	svc := syncsvc.New(client, database.DB)

	ctx := context.Background()
	date := time.Now().Format("2006-01-02")

	run := func(name string, fn func()) {
		if *fresh {
			_ = syncsvc.ClearCursor(ctx, name, date)
		}
		log.Printf("[backfill] starting %s (date=%s, interval=%s, fresh=%v)", name, date, *interval, *fresh)
		start := time.Now()
		fn()
		log.Printf("[backfill] finished %s in %s", name, time.Since(start).Round(time.Second))
	}

	switch *entity {
	case "institutes":
		run("institutes", func() { svc.SyncInstitutes(ctx, date, *interval) })
	case "staffs":
		run("staffs", func() { svc.SyncStaffs(ctx, date, *interval) })
	case "ranks":
		run("ranks", func() { svc.SyncRanks(ctx, date) })
	case "positions":
		run("positions", func() { svc.SyncPositions(ctx, date) })
	case "all":
		run("ranks", func() { svc.SyncRanks(ctx, date) })
		run("positions", func() { svc.SyncPositions(ctx, date) })
		run("institutes", func() { svc.SyncInstitutes(ctx, date, *interval) })
		run("staffs", func() { svc.SyncStaffs(ctx, date, *interval) })
	default:
		log.Fatalf("unknown -entity %q (use institutes|staffs|ranks|positions|all|staff-hierarchy)", *entity)
	}

	log.Println("[backfill] done")
}

// backfillStaffHierarchy recomputes institute_ids + institute_hierarchy for every
// existing staff. It groups by distinct institute_id so each chain is computed
// once and applied with a single UPDATE per institute_id (far fewer than 1 per row).
func backfillStaffHierarchy() {
	start := time.Now()
	ix, err := hierarchy.BuildIndex(database.DB)
	if err != nil {
		log.Fatalf("[backfill] build institute index failed: %v", err)
	}

	var instituteIDs []string
	if err := database.DB.Model(&models.Staff{}).
		Where("institute_id <> ''").
		Distinct().Pluck("institute_id", &instituteIDs).Error; err != nil {
		log.Fatalf("[backfill] list distinct institute_ids failed: %v", err)
	}
	log.Printf("[backfill] staff-hierarchy: %d distinct institute_ids to process", len(instituteIDs))

	var updated int64
	for i, iid := range instituteIDs {
		ids, names := ix.Chain(iid)
		res := database.DB.Model(&models.Staff{}).Where("institute_id = ?", iid).
			Updates(map[string]any{
				"institute_ids":       models.StringSlice(ids),
				"institute_hierarchy": models.StringSlice(names),
			})
		if res.Error != nil {
			log.Printf("[backfill] update for institute_id=%s failed: %v", iid, res.Error)
			continue
		}
		updated += res.RowsAffected
		if (i+1)%500 == 0 {
			log.Printf("[backfill] processed %d/%d institute_ids (%d staff rows updated)", i+1, len(instituteIDs), updated)
		}
	}
	log.Printf("[backfill] staff-hierarchy done: %d institute_ids, %d staff rows updated in %s",
		len(instituteIDs), updated, time.Since(start).Round(time.Second))
}

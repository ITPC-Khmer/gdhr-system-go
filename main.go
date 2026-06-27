package main

import (
	"context"
	"log"

	"backend/cache"
	"backend/config"
	"backend/database"
	"backend/handlers"
	"backend/routes"
	"backend/scheduler"
	"backend/services/gdhr"
	syncsvc "backend/services/sync"
)

func main() {
	cfg := config.Load()
	database.Connect(cfg)

	handlers.InitImageProxy(cfg.ImageProxyHosts)
	startSync(cfg)

	r := routes.Setup(cfg)

	addr := ":" + cfg.AppPort
	log.Printf("🚀 server running on http://localhost%s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}

// startSync wires up and starts the scheduled GDHR importer when enabled.
func startSync(cfg *config.Config) {
	if !cfg.SyncEnabled {
		log.Println("ℹ sync scheduler disabled (set SYNC_ENABLED=true to enable)")
		return
	}

	cache.Connect(cfg)
	client := gdhr.New(cfg.GDHRBaseURL, cfg.GDHRUser, cfg.GDHRPass, cfg.GDHRKey, cfg.SyncHTTPTimeout)
	svc := syncsvc.New(client, database.DB)

	sched, err := scheduler.New(cfg, svc)
	if err != nil {
		log.Fatalf("scheduler init failed: %v", err)
	}
	handlers.InitSync(sched, svc)
	sched.Start(context.Background())

	log.Printf("✓ sync scheduler started (window %s–%s %s)",
		cfg.SyncWindowStart, cfg.SyncWindowEnd, cfg.SyncTimezone)
}

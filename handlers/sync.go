package handlers

import (
	"context"
	"net/http"

	"backend/database"
	"backend/models"
	"backend/scheduler"
	syncsvc "backend/services/sync"

	"github.com/gin-gonic/gin"
)

// Sync dependencies are injected from main when SYNC_ENABLED=true. When nil,
// the trigger/status endpoints report that sync is disabled.
var (
	syncScheduler *scheduler.Scheduler
	syncService   *syncsvc.Service
)

// InitSync wires the scheduler/service into the handlers.
func InitSync(s *scheduler.Scheduler, svc *syncsvc.Service) {
	syncScheduler = s
	syncService = svc
}

// TriggerSync starts an on-demand sync run (admin only).
func TriggerSync(c *gin.Context) {
	if syncScheduler == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"message": "sync is disabled (set SYNC_ENABLED=true)"})
		return
	}
	started, day := syncScheduler.TriggerNow(context.Background())
	if !started {
		c.JSON(http.StatusConflict, gin.H{"message": "a sync run is already active", "date": day})
		return
	}
	c.JSON(http.StatusAccepted, gin.H{"message": "sync started", "date": day})
}

// SyncStatus reports row counts plus the Redis cursor state for today.
func SyncStatus(c *gin.Context) {
	var inst, staff, rank, pos int64
	database.DB.Model(&models.Institute{}).Count(&inst)
	database.DB.Model(&models.Staff{}).Count(&staff)
	database.DB.Model(&models.Rank{}).Count(&rank)
	database.DB.Model(&models.Position{}).Count(&pos)

	resp := gin.H{
		"enabled": syncScheduler != nil,
		"counts": gin.H{
			"institutes": inst,
			"staffs":     staff,
			"ranks":      rank,
			"positions":  pos,
		},
	}
	if syncScheduler != nil && syncService != nil {
		day := syncScheduler.Today()
		resp["date"] = day
		resp["cursors"] = syncService.CursorStatus(context.Background(), day)
	}
	c.JSON(http.StatusOK, resp)
}

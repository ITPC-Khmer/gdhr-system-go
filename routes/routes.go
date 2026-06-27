package routes

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"backend/config"
	"backend/handlers"
	"backend/middleware"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// Setup registers all routes and middleware, returning the gin engine.
func Setup(cfg *config.Config) *gin.Engine {
	if cfg.AppEnv != "development" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()

	// CORS
	corsCfg := cors.DefaultConfig()
	corsCfg.AllowOrigins = []string{cfg.CorsOrigin}
	corsCfg.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"}
	corsCfg.AllowHeaders = []string{"Origin", "Content-Type", "Authorization"}
	r.Use(cors.New(corsCfg))

	// Health check
	r.GET("/api/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	api := r.Group("/api")
	{
		// Public auth routes
		auth := api.Group("/auth")
		{
			auth.POST("/register", handlers.Register)
			auth.POST("/login", handlers.Login)
		}

		// Public image proxy (host-allowlisted) — usable directly from <img src>.
		api.GET("/image-proxy", handlers.ImageProxy)

		// Protected routes
		protected := api.Group("")
		protected.Use(middleware.AuthRequired())
		{
			protected.GET("/auth/me", handlers.Me)
			protected.GET("/stats", handlers.Stats)

			// User management is restricted to admins.
			users := protected.Group("/users")
			users.Use(middleware.AdminOnly())
			{
				users.GET("", handlers.ListUsers)
				users.GET("/:id", handlers.GetUser)
				users.POST("", handlers.CreateUser)
				users.PUT("/:id", handlers.UpdateUser)
				users.DELETE("/:id", handlers.DeleteUser)
			}

			// Synced GDHR data: list/get for any authenticated user; writes admin-only.
			registerCRUD(protected, "/institutes",
				handlers.ListInstitutes, handlers.GetInstitute,
				handlers.CreateInstitute, handlers.UpdateInstitute, handlers.DeleteInstitute)
			registerCRUD(protected, "/staffs",
				handlers.ListStaffs, handlers.GetStaff,
				handlers.CreateStaff, handlers.UpdateStaff, handlers.DeleteStaff)
			registerCRUD(protected, "/ranks",
				handlers.ListRanks, handlers.GetRank,
				handlers.CreateRank, handlers.UpdateRank, handlers.DeleteRank)
			registerCRUD(protected, "/positions",
				handlers.ListPositions, handlers.GetPosition,
				handlers.CreatePosition, handlers.UpdatePosition, handlers.DeletePosition)
			registerCRUD(protected, "/holidays",
				handlers.ListHolidays, handlers.GetHoliday,
				handlers.CreateHoliday, handlers.UpdateHoliday, handlers.DeleteHoliday)

			// Leave reference tables: list/get for any user, writes admin-only.
			registerCRUD(protected, "/leave-types",
				handlers.ListLeaveTypes, handlers.GetLeaveType,
				handlers.CreateLeaveType, handlers.UpdateLeaveType, handlers.DeleteLeaveType)
			registerCRUD(protected, "/leave-roles",
				handlers.ListLeaveRoles, handlers.GetLeaveRole,
				handlers.CreateLeaveRole, handlers.UpdateLeaveRole, handlers.DeleteLeaveRole)

			// Leave requests are read-only in-app (view list + approval
			// timeline). Any authenticated user can FILE a request (POST, which
			// auto-seeds the approval workflow); requests are never edited or
			// deleted through the API — the approval workflow changes their
			// state. Admins are view-only.
			leaves := protected.Group("/leaves")
			{
				leaves.GET("", handlers.ListLeaves)
				leaves.GET("/:id", handlers.GetLeave)
				leaves.POST("", handlers.CreateLeave)

				// Admin break-glass: force-resolve a STUCK request only.
				// Audited via approve_document; not the normal approval path.
				lw := leaves.Group("")
				lw.Use(middleware.AdminOnly())
				{
					lw.POST("/:id/override-approve", handlers.OverrideApproveLeave)
					lw.POST("/:id/override-reject", handlers.OverrideRejectLeave)
				}
			}

			registerCRUD(protected, "/leave-approvals",
				handlers.ListLeaveApprovals, handlers.GetLeaveApproval,
				handlers.CreateLeaveApproval, handlers.UpdateLeaveApproval, handlers.DeleteLeaveApproval)
			// Per-approver task actions: act on a single pending approval step.
			// An approver's task list is GET /leave-approvals?staff_id=&status=pending.
			protected.POST("/leave-approvals/:id/approve", handlers.ApproveLeaveApproval)
			protected.POST("/leave-approvals/:id/reject", handlers.RejectLeaveApproval)
			registerCRUD(protected, "/leave-files",
				handlers.ListLeaveFiles, handlers.GetLeaveFile,
				handlers.CreateLeaveFile, handlers.UpdateLeaveFile, handlers.DeleteLeaveFile)
			registerCRUD(protected, "/leave-years",
				handlers.ListLeaveYears, handlers.GetLeaveYear,
				handlers.CreateLeaveYear, handlers.UpdateLeaveYear, handlers.DeleteLeaveYear)

			// Staff<->institute management mapping (drives approval routing).
			registerCRUD(protected, "/staff-institute-roles",
				handlers.ListStaffInstituteRoles, handlers.GetStaffInstituteRole,
				handlers.CreateStaffInstituteRole, handlers.UpdateStaffInstituteRole, handlers.DeleteStaffInstituteRole)

			// Sync control is restricted to admins.
			syncGroup := protected.Group("/sync")
			syncGroup.Use(middleware.AdminOnly())
			{
				syncGroup.GET("/status", handlers.SyncStatus)
				syncGroup.POST("/trigger", handlers.TriggerSync)
			}
		}
	}

	// Serve the built Vue SPA from the same origin (single-domain setup).
	serveSPA(r, cfg.StaticDir)

	return r
}

// registerCRUD wires REST CRUD for a resource onto group g: list + get-one are
// available to any authenticated user; create/update/delete require admin.
func registerCRUD(g *gin.RouterGroup, path string, list, get, create, update, del gin.HandlerFunc) {
	res := g.Group(path)
	res.GET("", list)
	res.GET("/:id", get)

	write := res.Group("")
	write.Use(middleware.AdminOnly())
	{
		write.POST("", create)
		write.PUT("/:id", update)
		write.DELETE("/:id", del)
	}
}

// serveSPA serves static assets from staticDir and falls back to index.html
// for any non-API route so Vue Router (history mode) handles client routing.
func serveSPA(r *gin.Engine, staticDir string) {
	absDir, err := filepath.Abs(staticDir)
	if err != nil {
		absDir = staticDir
	}
	indexFile := filepath.Join(absDir, "index.html")

	r.NoRoute(func(c *gin.Context) {
		path := c.Request.URL.Path

		// Never let the SPA fallback swallow API calls.
		if strings.HasPrefix(path, "/api") {
			c.JSON(http.StatusNotFound, gin.H{"message": "route not found"})
			return
		}

		// Try to serve a real static file (assets, favicon, etc.).
		requested := filepath.Join(absDir, filepath.Clean("/"+path))
		if strings.HasPrefix(requested, absDir) {
			if info, statErr := os.Stat(requested); statErr == nil && !info.IsDir() {
				c.File(requested)
				return
			}
		}

		// Otherwise return index.html (or a hint if the frontend isn't built yet).
		if _, statErr := os.Stat(indexFile); statErr != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"message": "frontend not built — run `npm run build` in ./frontend (expected at " + indexFile + ")",
			})
			return
		}
		c.File(indexFile)
	})
}

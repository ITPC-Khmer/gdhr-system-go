package routes

import (
	"testing"

	"backend/config"
)

// TestSetupNoRouteConflicts ensures all routes register without gin panicking
// on conflicting wildcards (e.g. /leave-approvals/:id vs /:id/approve).
func TestSetupNoRouteConflicts(t *testing.T) {
	cfg := &config.Config{AppEnv: "development", CorsOrigin: "*", StaticDir: "."}
	if r := Setup(cfg); r == nil {
		t.Fatal("Setup returned nil engine")
	}
}

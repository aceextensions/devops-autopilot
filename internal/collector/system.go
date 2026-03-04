// Package collector wraps github.com/aceextensions/collector/system
// for use inside the devops-autopilot application.
package collector

import (
	"time"

	"github.com/aceextensions/collector/system"
)

// SystemStats re-exports the upstream type so the rest of the app
// doesn't need to import the external package directly.
type SystemStats = system.SystemStats

// CollectSystem collects live system metrics using the aceextensions/collector package.
func CollectSystem(serverName string) (*SystemStats, error) {
	return system.Collect(serverName)
}

// Ensure time is used (for the CollectedAt field reference elsewhere).
var _ = time.Now

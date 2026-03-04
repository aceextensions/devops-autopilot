// Package alert wraps github.com/aceextensions/alertengine/engine
// for use inside the devops-autopilot application.
package alert

import (
	"github.com/aceextensions/alertengine"
	"github.com/aceextensions/alertengine/engine"
	"github.com/aceextensions/devops-autopilot/internal/collector"
	"github.com/aceextensions/devops-autopilot/internal/config"
)

// Alert re-exports the upstream type.
type Alert = alertengine.Alert

// CheckThresholds builds Rules from live system stats + config thresholds
// and delegates evaluation to the aceextensions/alertengine package.
func CheckThresholds(cfg *config.Config, sys *collector.SystemStats) []Alert {
	rules := []alertengine.Rule{
		{Metric: "CPU", Value: sys.CPUPercent, Threshold: float64(cfg.Alerts.CPU)},
		{Metric: "Memory", Value: sys.MemPercent, Threshold: float64(cfg.Alerts.Memory)},
		{Metric: "Disk", Value: sys.DiskPercent, Threshold: float64(cfg.Alerts.Disk)},
	}
	return engine.Check(rules)
}

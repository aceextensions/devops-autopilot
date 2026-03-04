package report

import (
	"fmt"
	"strings"
	"time"

	"github.com/aceextensions/devops-autopilot/internal/alert"
	"github.com/aceextensions/devops-autopilot/internal/collector"
)

func Build(sys *collector.SystemStats, doc *collector.DockerStats, alerts []alert.Alert) string {
	var sb strings.Builder

	sb.WriteString("╔══════════════════════════════════════╗\n")
	sb.WriteString("║  DevOps Autopilot Report             ║\n")
	sb.WriteString("╚══════════════════════════════════════╝\n\n")

	sb.WriteString(fmt.Sprintf("Server   : %s\n", sys.ServerName))
	sb.WriteString(fmt.Sprintf("Time     : %s\n", time.Now().Format("2006-01-02 15:04:05")))
	sb.WriteString(fmt.Sprintf("Uptime   : %s\n\n", sys.Uptime))

	sb.WriteString("── System Health ─────────────────────\n")
	sb.WriteString(fmt.Sprintf("CPU      : %.1f%%\n", sys.CPUPercent))

	memUsedGB := float64(sys.MemUsed) / (1024 * 1024 * 1024)
	memTotalGB := float64(sys.MemTotal) / (1024 * 1024 * 1024)
	sb.WriteString(fmt.Sprintf("Memory   : %.1f%% (%.1fGB / %.1fGB)\n", sys.MemPercent, memUsedGB, memTotalGB))

	diskUsedGB := float64(sys.DiskUsed) / (1024 * 1024 * 1024)
	diskTotalGB := float64(sys.DiskTotal) / (1024 * 1024 * 1024)
	sb.WriteString(fmt.Sprintf("Disk     : %.1f%% (%.1fGB / %.1fGB)\n", sys.DiskPercent, diskUsedGB, diskTotalGB))
	sb.WriteString(fmt.Sprintf("Load Avg : %.2f / %.2f / %.2f\n\n", sys.LoadAvg1, sys.LoadAvg5, sys.LoadAvg15))

	sb.WriteString("── Docker ────────────────────────────\n")
	if !doc.DaemonAvailable {
		sb.WriteString("Status   : Not Available\n\n")
	} else {
		sb.WriteString(fmt.Sprintf("Running  : %d\n", doc.Running))
		sb.WriteString(fmt.Sprintf("Stopped  : %d\n", doc.Stopped))
		sb.WriteString(fmt.Sprintf("Total    : %d\n\n", doc.Total))
		for _, c := range doc.Containers {
			icon := "○"
			if c.State == "running" {
				icon = "●"
			}
			sb.WriteString(fmt.Sprintf("  %s [%s] %s (%s) restarts=%d\n", icon, c.ID, c.Name, c.State, c.RestartCount))
		}
		sb.WriteString("\n")
	}

	sb.WriteString("── Alerts ─────────────────────────────\n")
	if len(alerts) == 0 {
		sb.WriteString("None\n")
	} else {
		for _, a := range alerts {
			sb.WriteString(fmt.Sprintf("⚠ %s\n", a.Message))
		}
	}

	return sb.String()
}

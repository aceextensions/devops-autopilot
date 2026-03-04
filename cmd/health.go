package cmd

import (
	"fmt"

	"github.com/aceextensions/devops-autopilot/internal/collector"
	"github.com/aceextensions/devops-autopilot/internal/config"
	"github.com/spf13/cobra"
)

var healthCmd = &cobra.Command{
	Use:   "health",
	Short: "Show current system health metrics",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load(cfgFile)
		if err != nil {
			return err
		}

		stats, err := collector.CollectSystem(cfg.ServerName)
		if err != nil {
			return err
		}

		fmt.Printf("Server   : %s\n", stats.ServerName)
		fmt.Printf("Uptime   : %s\n", stats.Uptime)
		fmt.Printf("CPU      : %.1f%%\n", stats.CPUPercent)
		memUsedGB := float64(stats.MemUsed) / (1024 * 1024 * 1024)
		memTotalGB := float64(stats.MemTotal) / (1024 * 1024 * 1024)
		fmt.Printf("Memory   : %.1f%% (used %.1fGB / total %.1fGB)\n", stats.MemPercent, memUsedGB, memTotalGB)
		diskUsedGB := float64(stats.DiskUsed) / (1024 * 1024 * 1024)
		diskTotalGB := float64(stats.DiskTotal) / (1024 * 1024 * 1024)
		fmt.Printf("Disk     : %.1f%% (used %.1fGB / total %.1fGB)\n", stats.DiskPercent, diskUsedGB, diskTotalGB)
		fmt.Printf("Load Avg : %.2f / %.2f / %.2f\n", stats.LoadAvg1, stats.LoadAvg5, stats.LoadAvg15)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(healthCmd)
}

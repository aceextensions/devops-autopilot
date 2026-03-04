package cmd

import (
	"fmt"

	"github.com/aceextensions/devops-autopilot/internal/alert"
	"github.com/aceextensions/devops-autopilot/internal/collector"
	"github.com/aceextensions/devops-autopilot/internal/config"
	"github.com/aceextensions/devops-autopilot/internal/report"
	"github.com/spf13/cobra"
)

var reportCmd = &cobra.Command{
	Use:   "report",
	Short: "Generate a full server health report",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load(cfgFile)
		if err != nil {
			return err
		}

		sys, err := collector.CollectSystem(cfg.ServerName)
		if err != nil {
			return err
		}

		doc, _ := collector.CollectDocker()
		alerts := alert.CheckThresholds(cfg, sys)

		rpt := report.Build(sys, doc, alerts)
		fmt.Println(rpt)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(reportCmd)
}

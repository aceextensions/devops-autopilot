package cmd

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/aceextensions/devops-autopilot/internal/alert"
	"github.com/aceextensions/devops-autopilot/internal/collector"
	"github.com/aceextensions/devops-autopilot/internal/config"
	"github.com/aceextensions/devops-autopilot/internal/notifier"
	"github.com/aceextensions/devops-autopilot/internal/report"
	"github.com/spf13/cobra"
)

var monitorCmd = &cobra.Command{
	Use:   "monitor",
	Short: "Run continuous monitoring loop",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load(cfgFile)
		if err != nil {
			return err
		}

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

		go func() {
			<-sigCh
			fmt.Println("\nShutting down monitor...")
			cancel()
		}()

		fmt.Printf("Starting monitor (interval: %ds)...\n", cfg.Monitor.IntervalSeconds)
		ticker := time.NewTicker(time.Duration(cfg.Monitor.IntervalSeconds) * time.Second)
		defer ticker.Stop()

		slackNs := notifier.NewSlack(cfg)
		emailNs := notifier.NewEmail(cfg)

		runCheck(cfg, slackNs, emailNs)

		for {
			select {
			case <-ctx.Done():
				return nil
			case <-ticker.C:
				runCheck(cfg, slackNs, emailNs)
			}
		}
	},
}

func runCheck(cfg *config.Config, slackNs *notifier.SlackNotifier, emailNs *notifier.EmailNotifier) {
	sys, err := collector.CollectSystem(cfg.ServerName)
	if err != nil {
		fmt.Printf("Error collecting system stats: %v\n", err)
		return
	}

	doc, _ := collector.CollectDocker()
	alerts := alert.CheckThresholds(cfg, sys)

	if len(alerts) > 0 {
		rpt := report.Build(sys, doc, alerts)
		_ = slackNs.Send(rpt)
		_ = emailNs.Send(fmt.Sprintf("Alerts on %s", cfg.ServerName), rpt)
		fmt.Printf("[%s] %d alerts triggered\n", time.Now().Format("15:04:05"), len(alerts))
	} else {
		fmt.Printf("[%s] System healthy\n", time.Now().Format("15:04:05"))
	}
}

func init() {
	rootCmd.AddCommand(monitorCmd)
}

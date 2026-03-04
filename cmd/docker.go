package cmd

import (
	"fmt"

	"github.com/aceextensions/devops-autopilot/internal/collector"
	"github.com/spf13/cobra"
)

var dockerCmd = &cobra.Command{
	Use:   "docker",
	Short: "Show Docker container status",
	RunE: func(cmd *cobra.Command, args []string) error {
		stats, err := collector.CollectDocker()
		if err != nil {
			return err
		}

		if !stats.DaemonAvailable {
			fmt.Println("Docker Status  : Not Available")
			return nil
		}

		fmt.Println("Docker Status  : Available")
		fmt.Printf("Running        : %d\n", stats.Running)
		fmt.Printf("Stopped        : %d\n", stats.Stopped)
		fmt.Printf("Total          : %d\n", stats.Total)
		fmt.Println()

		for _, c := range stats.Containers {
			icon := "○"
			if c.State == "running" {
				icon = "●"
			}
			fmt.Printf("  %s [%s] %s (%s) restarts=%d\n", icon, c.ID, c.Name, c.State, c.RestartCount)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(dockerCmd)
}

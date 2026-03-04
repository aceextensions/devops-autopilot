package cmd

import (
	"fmt"
	"os"

	"github.com/aceextensions/devops-autopilot/internal/config"
	svcmanager "github.com/aceextensions/devops-autopilot/internal/service"
	"github.com/spf13/cobra"
)

var serviceCmd = &cobra.Command{
	Use:   "service [install|uninstall|start|stop|restart|status|run]",
	Short: "Manage the system service",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		action := args[0]

		// Load config to just make sure path is valid before installing/running
		_, err := config.Load(cfgFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Warning: failed to load config at %s: %v\n", cfgFile, err)
		}

		err = svcmanager.Manage(action, cfgFile)
		if err != nil {
			return fmt.Errorf("service error: %w", err)
		}

		if action != "run" && action != "status" {
			fmt.Printf("Service %s complete.\n", action)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(serviceCmd)
}

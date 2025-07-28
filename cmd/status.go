package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"stackpulse/internal/monitor"
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show current resource usage of monitored services",
	Long:  `Display real-time resource usage statistics for all monitored Node.js services.`,
	RunE:  runStatus,
}

func init() {
	rootCmd.AddCommand(statusCmd)
}

func runStatus(cmd *cobra.Command, args []string) error {
	status, err := monitor.GetCurrentStatus()
	if err != nil {
		return fmt.Errorf("failed to get status: %w", err)
	}

	monitor.DisplayStatus(status)
	return nil
}
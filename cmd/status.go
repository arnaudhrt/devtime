package cmd

import (
	"github.com/arnaudhrt/devtime/internal"
	"time"

	"github.com/spf13/cobra"
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show current coding status",
	RunE: func(cmd *cobra.Command, args []string) error {
		now := time.Now()
		start := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
		events, err := internal.ReadEventsForRange(start, now)
		if err != nil {
			return err
		}

		sessions := internal.ComputeSessions(events)
		if len(sessions) == 0 {
			internal.PrintStatus(false, nil)
			return nil
		}

		last := &sessions[len(sessions)-1]
		active := time.Since(last.End) < 5*time.Minute
		internal.PrintStatus(active, last)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(statusCmd)
}

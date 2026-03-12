package cmd

import (
	"github.com/arnaudhrt/devtime/internal"

	"github.com/spf13/cobra"
)

var todayCmd = &cobra.Command{
	Use:   "today",
	Short: "Show today's coding time",
	RunE: func(cmd *cobra.Command, args []string) error {
		start, end := internal.TodayRange()
		events, err := internal.ReadEventsForRange(start, end)
		if err != nil {
			return err
		}
		sessions := internal.ComputeSessions(events)
		summary := internal.Summarize(sessions)
		internal.PrintSummary("Today", summary)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(todayCmd)
}

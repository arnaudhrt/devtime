package cmd

import (
	"github.com/arnaudhrt/devtime/internal"

	"github.com/spf13/cobra"
)

var weekCmd = &cobra.Command{
	Use:   "week",
	Short: "Show this week's coding time",
	RunE: func(cmd *cobra.Command, args []string) error {
		start, end := internal.WeekRange()
		events, err := internal.ReadEventsForRange(start, end)
		if err != nil {
			return err
		}
		sessions := internal.ComputeSessions(events)
		summary := internal.Summarize(sessions)
		internal.PrintSummary("This Week", summary)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(weekCmd)
}

package cmd

import (
	"github.com/arnaudhrt/devtime/internal"

	"github.com/spf13/cobra"
)

var monthCmd = &cobra.Command{
	Use:   "month",
	Short: "Show this month's coding time",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := internal.CheckDataExists(); err != nil {
			return err
		}
		start, end := internal.MonthRange()
		events, err := internal.ReadEventsForRange(start, end)
		if err != nil {
			return err
		}
		sessions := internal.ComputeSessions(events)
		summary := internal.Summarize(sessions)
		internal.PrintSummary("This Month", summary)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(monthCmd)
}

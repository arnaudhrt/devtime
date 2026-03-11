package cmd

import (
	"devtime/internal"
	"fmt"

	"github.com/spf13/cobra"
)

var projectCmd = &cobra.Command{
	Use:   "project <name> <all|month|week>",
	Short: "Show coding time for a specific project",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]
		period := args[1]

		events, header, err := loadEventsForPeriod(period)
		if err != nil {
			return err
		}

		sessions := internal.ComputeSessions(events)
		sessions = internal.FilterByProject(sessions, name)
		summary := internal.Summarize(sessions)
		internal.PrintSummary(fmt.Sprintf("%s — %s", name, header), summary)
		return nil
	},
}

func loadEventsForPeriod(period string) ([]internal.Event, string, error) {
	switch period {
	case "week":
		start, end := internal.WeekRange()
		events, err := internal.ReadEventsForRange(start, end)
		return events, "This Week", err
	case "month":
		start, end := internal.MonthRange()
		events, err := internal.ReadEventsForRange(start, end)
		return events, "This Month", err
	case "all":
		events, err := internal.ReadAllEvents()
		return events, "All Time", err
	default:
		return nil, "", fmt.Errorf("invalid period %q: use all, month, or week", period)
	}
}

func init() {
	rootCmd.AddCommand(projectCmd)
}

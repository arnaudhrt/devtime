package cmd

import (
	"fmt"
	"strings"

	"github.com/arnaudhrt/devtime/internal"
	"github.com/spf13/cobra"
)

var projectCmd = &cobra.Command{
	Use:   "project <name>",
	Short: "Show coding time for a specific project",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := internal.CheckDataExists(); err != nil {
			return err
		}
		name := args[0]

		// All time
		allEvents, err := internal.ReadAllEvents()
		if err != nil {
			return err
		}
		allSessions := internal.FilterByProject(internal.ComputeSessions(allEvents), name)
		allSummary := internal.Summarize(allSessions)

		if allSummary.Total == 0 {
			fmt.Printf("No data for project %q.\n", name)
			return nil
		}

		// This month
		mStart, mEnd := internal.MonthRange()
		monthEvents, err := internal.ReadEventsForRange(mStart, mEnd)
		if err != nil {
			return err
		}
		monthSummary := internal.Summarize(internal.FilterByProject(internal.ComputeSessions(monthEvents), name))

		// This week
		wStart, wEnd := internal.WeekRange()
		weekEvents, err := internal.ReadEventsForRange(wStart, wEnd)
		if err != nil {
			return err
		}
		weekSummary := internal.Summarize(internal.FilterByProject(internal.ComputeSessions(weekEvents), name))

		// Header
		fmt.Printf("\n  Devtime for %s\n\n", name)

		// Period breakdown
		fmt.Printf("  All time:    %s\n", internal.FormatDuration(allSummary.Total))
		fmt.Printf("  This month:  %s\n", internal.FormatDuration(monthSummary.Total))
		fmt.Printf("  This week:   %s\n", internal.FormatDuration(weekSummary.Total))

		// Languages list with bars
		if len(allSummary.Languages) > 0 {
			maxLen := 0
			for _, l := range allSummary.Languages {
				if len(l.Name) > maxLen {
					maxLen = len(l.Name)
				}
			}

			fmt.Println()
			fmt.Println("  Languages:")
			for _, l := range allSummary.Languages {
				pct := float64(l.Duration) / float64(allSummary.Total) * 100
				filled := int(float64(internal.BarWidth) * float64(l.Duration) / float64(allSummary.Total))
				if filled < 0 {
					filled = 0
				}
				if filled > internal.BarWidth {
					filled = internal.BarWidth
				}
				empty := internal.BarWidth - filled
				bar := strings.Repeat("█", filled) + strings.Repeat("░", empty)
				fmt.Printf("    %-*s  %s  %s  %3.0f%%\n", maxLen, l.Name, internal.FormatDuration(l.Duration), bar, pct)
			}
		}

		fmt.Println()
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

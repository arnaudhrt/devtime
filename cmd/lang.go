package cmd

import (
	"fmt"
	"strings"

	"github.com/arnaudhrt/devtime/internal"
	"github.com/spf13/cobra"
)

var langCmd = &cobra.Command{
	Use:   "lang <name>",
	Short: "Show coding time for a specific language",
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
		allSessions := internal.FilterByLanguage(internal.ComputeSessions(allEvents), name)
		allSummary := internal.Summarize(allSessions)

		if allSummary.Total == 0 {
			fmt.Printf("No data for language %q.\n", name)
			return nil
		}

		// This month
		mStart, mEnd := internal.MonthRange()
		monthEvents, err := internal.ReadEventsForRange(mStart, mEnd)
		if err != nil {
			return err
		}
		monthSummary := internal.Summarize(internal.FilterByLanguage(internal.ComputeSessions(monthEvents), name))

		// This week
		wStart, wEnd := internal.WeekRange()
		weekEvents, err := internal.ReadEventsForRange(wStart, wEnd)
		if err != nil {
			return err
		}
		weekSummary := internal.Summarize(internal.FilterByLanguage(internal.ComputeSessions(weekEvents), name))

		// Header
		fmt.Printf("\n  Devtime for %s\n\n", strings.ToUpper(name))

		// Period breakdown
		fmt.Printf("  All time:    %s\n", internal.FormatDuration(allSummary.Total))
		fmt.Printf("  This month:  %s\n", internal.FormatDuration(monthSummary.Total))
		fmt.Printf("  This week:   %s\n", internal.FormatDuration(weekSummary.Total))

		// Projects list
		if len(allSummary.Projects) > 0 {
			maxLen := 0
			for _, p := range allSummary.Projects {
				if len(p.Name) > maxLen {
					maxLen = len(p.Name)
				}
			}

			fmt.Println()
			fmt.Println("  Projects:")
			for _, p := range allSummary.Projects {
				fmt.Printf("    %-*s  %s\n", maxLen, p.Name, internal.FormatDuration(p.Duration))
			}
		}

		fmt.Println()
		return nil
	},
}

func init() {
	rootCmd.AddCommand(langCmd)
}

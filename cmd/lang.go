package cmd

import (
	"fmt"

	"github.com/arnaudhrt/devtime/internal"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

var langCmd = &cobra.Command{
	Use:   "lang [name]",
	Short: "Show coding time for a language (interactive if no name given)",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := internal.CheckDataExists(); err != nil {
			return err
		}

		if len(args) == 0 {
			langs, err := internal.AllTimeLanguageNames()
			if err != nil {
				return err
			}
			if len(langs) == 0 {
				fmt.Println("No languages found.")
				return nil
			}

			prompt := promptui.Select{
				Label: "Select a language",
				Items: langs,
				Size:  15,
			}

			_, selected, err := prompt.Run()
			if err != nil {
				fmt.Println("Selection cancelled.")
				return nil
			}
			args = []string{selected}
		}

		name := args[0]

		// All time
		allSummary, err := internal.AllTimeSummaryForLanguage(name)
		if err != nil {
			return err
		}

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

		// Period breakdown
		fmt.Printf("  \n")
		fmt.Printf("  All time:    %s\n", internal.FormatDuration(allSummary.Total))
		fmt.Printf("  This month:  %s\n", internal.FormatDuration(monthSummary.Total))
		fmt.Printf("  This week:   %s\n", internal.FormatDuration(weekSummary.Total))

		allSummary = internal.FilterShort(allSummary)

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

package cmd

import (
	"fmt"
	"time"

	"github.com/arnaudhrt/devtime/internal"
	"github.com/spf13/cobra"
)

var allCmd = &cobra.Command{
	Use:   "all",
	Short: "Show your overall coding summary",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := internal.CheckDataExists(); err != nil {
			return err
		}

		data, err := internal.LoadProfileData()
		if err != nil {
			return err
		}

		summary := data.Summary
		if summary.Total == 0 {
			fmt.Println("No coding data yet.")
			return nil
		}

		dailyAvg := summary.Total / time.Duration(data.DaysTracked)

		fmt.Printf("  \n")
		fmt.Printf("  Tracking since: %s\n", data.FirstDay.Format("Jan 2, 2006"))
		fmt.Printf("  Total time:     %s\n", internal.FormatDuration(summary.Total))
		fmt.Printf("  Daily average:  %s\n", internal.FormatDuration(dailyAvg))
		fmt.Printf("  Days tracked:   %d\n", data.DaysTracked)

		// Top 3 Projects
		if len(summary.Projects) > 0 {
			projects := summary.Projects
			if len(projects) > 3 {
				projects = projects[:3]
			}
			maxLen := 0
			for _, p := range projects {
				if len(p.Name) > maxLen {
					maxLen = len(p.Name)
				}
			}

			fmt.Println()
			fmt.Println("  Top 3 projects:")
			for _, p := range projects {
				fmt.Printf("    %-*s  %s\n", maxLen, p.Name, internal.FormatDuration(p.Duration))
			}
		}

		// Top 3 Languages
		if len(summary.Languages) > 0 {
			languages := summary.Languages
			if len(languages) > 3 {
				languages = languages[:3]
			}
			maxLen := 0
			for _, l := range languages {
				if len(l.Name) > maxLen {
					maxLen = len(l.Name)
				}
			}

			fmt.Println()
			fmt.Println("  Top 3 languages:")
			for _, l := range languages {
				fmt.Printf("    %-*s  %s\n", maxLen, l.Name, internal.FormatDuration(l.Duration))
			}
		}

		fmt.Println()
		return nil
	},
}

func init() {
	rootCmd.AddCommand(allCmd)
}

package cmd

import (
	"fmt"
	"strings"
	"time"

	"github.com/arnaudhrt/devtime/internal"
	"github.com/spf13/cobra"
)

var profileCmd = &cobra.Command{
	Use:   "profile",
	Short: "Show your overall coding profile",
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

		// Header
		fmt.Printf("  \nTotal time:     %s\n", internal.FormatDuration(summary.Total))
		fmt.Printf("  Daily average:  %s\n", internal.FormatDuration(dailyAvg))
		fmt.Printf("  Days tracked:   %d\n", data.DaysTracked)
		fmt.Printf("  Tracking since: %s\n", data.FirstDay.Format("Jan 2, 2006"))

		// Projects
		if len(summary.Projects) > 0 {
			maxLen := 0
			for _, p := range summary.Projects {
				if len(p.Name) > maxLen {
					maxLen = len(p.Name)
				}
			}

			fmt.Println()
			fmt.Println("  Projects:")
			for _, p := range summary.Projects {
				pct := float64(p.Duration) / float64(summary.Total) * 100
				filled := int(float64(internal.BarWidth) * float64(p.Duration) / float64(summary.Total))
				if filled < 0 {
					filled = 0
				}
				if filled > internal.BarWidth {
					filled = internal.BarWidth
				}
				empty := internal.BarWidth - filled
				bar := strings.Repeat("█", filled) + strings.Repeat("░", empty)
				fmt.Printf("    %-*s  %s  %s  %3.0f%%\n", maxLen, p.Name, internal.FormatDuration(p.Duration), bar, pct)
			}
		}

		// Languages
		if len(summary.Languages) > 0 {
			maxLen := 0
			for _, l := range summary.Languages {
				if len(l.Name) > maxLen {
					maxLen = len(l.Name)
				}
			}

			fmt.Println()
			fmt.Println("  Languages:")
			for _, l := range summary.Languages {
				pct := float64(l.Duration) / float64(summary.Total) * 100
				filled := int(float64(internal.BarWidth) * float64(l.Duration) / float64(summary.Total))
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

func init() {
	rootCmd.AddCommand(profileCmd)
}

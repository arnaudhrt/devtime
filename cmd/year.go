package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/arnaudhrt/devtime/internal"

	"github.com/spf13/cobra"
)

var yearCmd = &cobra.Command{
	Use:   "year [yyyy]",
	Short: "Show a year's coding time (default: current year)",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := internal.CheckDataExists(); err != nil {
			return err
		}

		now := time.Now()
		year := now.Year()
		header := "This Year"
		var lastMonth time.Month
		var end time.Time

		if len(args) == 1 {
			y, err := strconv.Atoi(args[0])
			if err != nil {
				return fmt.Errorf("invalid year: %q", args[0])
			}
			year = y
			header = fmt.Sprintf("%d", year)
		}

		if year == now.Year() {
			lastMonth = now.Month()
			end = now
		} else {
			lastMonth = time.December
			end = time.Date(year, time.December, 31, 23, 59, 59, 999999999, now.Location())
		}

		dir, err := internal.EventsDir()
		if err != nil {
			return err
		}

		combined := internal.Summary{}

		for m := time.January; m <= lastMonth; m++ {
			monthStr := fmt.Sprintf("%04d-%02d", year, int(m))
			summaryPath := filepath.Join(dir, fmt.Sprintf("summary-%s.json", monthStr))

			if _, err := os.Stat(summaryPath); err == nil {
				ms, err := internal.ReadMonthlySummary(summaryPath)
				if err != nil {
					return err
				}
				combined = internal.MergeSummary(combined, internal.SummaryFromMonthly(ms))
			}

			// Also read raw events for this month (both sources may exist).
			start := time.Date(year, m, 1, 0, 0, 0, 0, now.Location())
			monthEnd := start.AddDate(0, 1, 0).Add(-time.Nanosecond)
			if m == lastMonth && year == now.Year() {
				monthEnd = end
			}
			events, err := internal.ReadEventsForRange(start, monthEnd)
			if err != nil {
				return err
			}
			if len(events) > 0 {
				sessions := internal.ComputeSessions(events)
				combined = internal.MergeSummary(combined, internal.Summarize(sessions))
			}
		}

		internal.PrintSummary(header, combined)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(yearCmd)
}

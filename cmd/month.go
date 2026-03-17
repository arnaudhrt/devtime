package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/arnaudhrt/devtime/internal"
	"github.com/spf13/cobra"
)

var monthCmd = &cobra.Command{
	Use:   "month [mmm-yyyy]",
	Short: "Show a month's coding time (default: current month)",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := internal.CheckDataExists(); err != nil {
			return err
		}

		now := time.Now()
		var year int
		var month time.Month
		var header string

		if len(args) == 0 {
			year = now.Year()
			month = now.Month()
			header = "This Month"
		} else {
			// Parse mmm-yyyy (e.g. nov-2025)
			parts := strings.SplitN(strings.ToLower(args[0]), "-", 2)
			if len(parts) != 2 {
				return fmt.Errorf("invalid format: use mmm-yyyy (e.g. nov-2025)")
			}

			monthAbbr := strings.ToUpper(parts[0][:1]) + parts[0][1:]
			parsed, err := time.Parse("Jan", monthAbbr)
			if err != nil {
				return fmt.Errorf("invalid month: %q (use 3-letter abbreviation like jan, feb, mar)", parts[0])
			}
			month = parsed.Month()

			y, err := strconv.Atoi(parts[1])
			if err != nil {
				return fmt.Errorf("invalid year: %q", parts[1])
			}
			year = y
			header = fmt.Sprintf("%s %d", month.String(), year)
		}

		dir, err := internal.EventsDir()
		if err != nil {
			return err
		}

		combined := internal.Summary{}

		// Read summary if it exists.
		monthStr := fmt.Sprintf("%04d-%02d", year, int(month))
		summaryPath := filepath.Join(dir, fmt.Sprintf("summary-%s.json", monthStr))
		if _, err := os.Stat(summaryPath); err == nil {
			ms, err := internal.ReadMonthlySummary(summaryPath)
			if err != nil {
				return err
			}
			combined = internal.MergeSummary(combined, internal.SummaryFromMonthly(ms))
		}

		// Read raw events if they exist.
		start := time.Date(year, month, 1, 0, 0, 0, 0, now.Location())
		end := start.AddDate(0, 1, 0).Add(-time.Nanosecond)
		if year == now.Year() && month == now.Month() {
			end = now
		}
		events, err := internal.ReadEventsForRange(start, end)
		if err != nil {
			return err
		}
		if len(events) > 0 {
			sessions := internal.ComputeSessions(events)
			combined = internal.MergeSummary(combined, internal.Summarize(sessions))
		}

		internal.PrintSummary(header, combined)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(monthCmd)
}

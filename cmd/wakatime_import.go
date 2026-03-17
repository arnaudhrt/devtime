package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/arnaudhrt/devtime/internal"
	"github.com/spf13/cobra"
)

var wakatimeImportCmd = &cobra.Command{
	Use:   "wakatime-import <file>",
	Short: "Import coding history from a WakaTime JSON export",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		filePath := args[0]

		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			return fmt.Errorf("file not found: %s", filePath)
		}

		fmt.Printf("Parsing WakaTime export: %s\n", filePath)

		summaries, err := internal.ParseWakaTimeExport(filePath)
		if err != nil {
			return err
		}

		dir, err := internal.EventsDir()
		if err != nil {
			return err
		}

		// Ensure ~/.devtime/ exists.
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("creating data directory: %w", err)
		}

		now := time.Now()
		currentMonth := fmt.Sprintf("%04d-%02d", now.Year(), int(now.Month()))

		var totalSeconds int64
		var totalMonths int

		for _, imported := range summaries {
			var year, month int
			if _, err := fmt.Sscanf(imported.Month, "%d-%d", &year, &month); err == nil {
				eventsPath := internal.EventFilePath(dir, year, time.Month(month))
				if _, err := os.Stat(eventsPath); err == nil {
					if imported.Month == currentMonth {
						// Never compact the current month — raw events are
						// needed for today/week queries. The WakaTime data
						// goes into the summary and both sources are merged
						// at query time.
						fmt.Printf("  %s: skipping compaction (current month)\n", imported.Month)
					} else {
						fmt.Printf("  Compacting raw events for %s...\n", imported.Month)
						if err := internal.CompactMonth(dir, year, time.Month(month)); err != nil {
							return fmt.Errorf("compacting %s: %w", imported.Month, err)
						}
					}
				}
			}

			// Check if summary already exists — merge if so.
			summaryPath := filepath.Join(dir, fmt.Sprintf("summary-%s.json", imported.Month))
			if _, err := os.Stat(summaryPath); err == nil {
				existing, err := internal.ReadMonthlySummary(summaryPath)
				if err != nil {
					return fmt.Errorf("reading existing summary %s: %w", imported.Month, err)
				}
				imported = internal.MergeMonthlySummaries(existing, imported)
				fmt.Printf("  %s: merged with existing data → %s\n", imported.Month, formatSeconds(imported.TotalSeconds))
			} else {
				fmt.Printf("  %s: %s\n", imported.Month, formatSeconds(imported.TotalSeconds))
			}

			if err := internal.WriteSummary(dir, imported); err != nil {
				return fmt.Errorf("writing summary %s: %w", imported.Month, err)
			}

			totalSeconds += imported.TotalSeconds
			totalMonths++
		}

		fmt.Printf("\nImported %d month(s), total: %s\n", totalMonths, formatSeconds(totalSeconds))
		return nil
	},
}

func init() {
	rootCmd.AddCommand(wakatimeImportCmd)
}

func formatSeconds(secs int64) string {
	h := secs / 3600
	m := (secs % 3600) / 60
	return fmt.Sprintf("%dh %02dm", h, m)
}

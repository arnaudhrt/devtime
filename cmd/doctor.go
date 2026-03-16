package cmd

import (
	"fmt"
	"time"

	"github.com/arnaudhrt/devtime/internal"

	"github.com/spf13/cobra"
)

var doctorCmd = &cobra.Command{
	Use:   "doctor",
	Short: "Check if the VS Code extension is working",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := internal.CheckDataExists(); err != nil {
			return err
		}

		now := time.Now()
		start := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
		events, err := internal.ReadEventsForRange(start, now)
		if err != nil {
			return err
		}

		if len(events) == 0 {
			fmt.Println("\n  No events found this month.")
			fmt.Println("  Make sure the devtime VS Code extension is installed and running.")
			fmt.Println()
			return nil
		}

		last := events[len(events)-1]
		ago := time.Since(last.Timestamp)

		fmt.Println()
		fmt.Printf("  Last event:\n")
		fmt.Printf("    Type:     %s\n", last.Type)
		fmt.Printf("    Project:  %s\n", last.Project)
		fmt.Printf("    Language: %s\n", last.Language)
		fmt.Printf("    Editor:   %s\n", last.Editor)
		fmt.Printf("    Time:     %s (%s ago)\n", last.Timestamp.Format("Jan 02 15:04:05"), formatAgo(ago))
		fmt.Println()

		return nil
	},
}

func formatAgo(d time.Duration) string {
	if d < time.Minute {
		return fmt.Sprintf("%ds", int(d.Seconds()))
	}
	if d < time.Hour {
		return fmt.Sprintf("%dm", int(d.Minutes()))
	}
	if d < 24*time.Hour {
		h := int(d.Hours())
		m := int(d.Minutes()) % 60
		if m == 0 {
			return fmt.Sprintf("%dh", h)
		}
		return fmt.Sprintf("%dh %dm", h, m)
	}
	days := int(d.Hours()) / 24
	return fmt.Sprintf("%dd", days)
}

func init() {
	rootCmd.AddCommand(doctorCmd)
}

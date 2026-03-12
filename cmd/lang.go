package cmd

import (
	"github.com/arnaudhrt/devtime/internal"
	"fmt"

	"github.com/spf13/cobra"
)

var langCmd = &cobra.Command{
	Use:   "lang <name> <all|month|week>",
	Short: "Show coding time for a specific language",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]
		period := args[1]

		events, header, err := loadEventsForPeriod(period)
		if err != nil {
			return err
		}

		sessions := internal.ComputeSessions(events)
		sessions = internal.FilterByLanguage(sessions, name)
		summary := internal.Summarize(sessions)
		internal.PrintSummary(fmt.Sprintf("%s — %s", name, header), summary)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(langCmd)
}

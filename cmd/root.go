package cmd

import (
	"fmt"
	"os"

	"github.com/arnaudhrt/devtime/internal"
	"github.com/spf13/cobra"
)

var Version = "dev"

var rootCmd = &cobra.Command{
	Use:   "devtime",
	Short: "Track your coding time from the terminal",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := internal.CheckDataExists(); err != nil {
			return err
		}
		start, end := internal.TodayRange()
		events, err := internal.ReadEventsForRange(start, end)
		if err != nil {
			return err
		}
		sessions := internal.ComputeSessions(events)
		summary := internal.Summarize(sessions)
		internal.PrintSummary("Today", summary)
		return nil
	},
}

func init() {
	rootCmd.Version = Version
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

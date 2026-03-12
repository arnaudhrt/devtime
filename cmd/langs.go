package cmd

import (
	"fmt"
	"sort"

	"github.com/arnaudhrt/devtime/internal"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

var langsCmd = &cobra.Command{
	Use:   "langs",
	Short: "Interactively select a language and show its time summary",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := internal.CheckDataExists(); err != nil {
			return err
		}
		events, err := internal.ReadAllEvents()
		if err != nil {
			return err
		}

		sessions := internal.ComputeSessions(events)

		// Extract unique language names.
		seen := make(map[string]bool)
		for _, s := range sessions {
			seen[s.Language] = true
		}
		if len(seen) == 0 {
			fmt.Println("No languages found.")
			return nil
		}

		langs := make([]string, 0, len(seen))
		for name := range seen {
			langs = append(langs, name)
		}
		sort.Strings(langs)

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

		return langCmd.RunE(cmd, []string{selected})
	},
}

func init() {
	rootCmd.AddCommand(langsCmd)
}

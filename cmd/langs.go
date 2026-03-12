package cmd

import (
	"fmt"

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

		return langCmd.RunE(cmd, []string{selected})
	},
}

func init() {
	rootCmd.AddCommand(langsCmd)
}

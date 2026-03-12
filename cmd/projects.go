package cmd

import (
	"fmt"

	"github.com/arnaudhrt/devtime/internal"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

var projectsCmd = &cobra.Command{
	Use:   "projects",
	Short: "Interactively select a project and show its time summary",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := internal.CheckDataExists(); err != nil {
			return err
		}
		projects, err := internal.AllTimeProjectNames()
		if err != nil {
			return err
		}
		if len(projects) == 0 {
			fmt.Println("No projects found.")
			return nil
		}

		prompt := promptui.Select{
			Label: "Select a project",
			Items: projects,
			Size:  15,
		}

		_, selected, err := prompt.Run()
		if err != nil {
			fmt.Println("Selection cancelled.")
			return nil
		}

		return projectCmd.RunE(cmd, []string{selected})
	},
}

func init() {
	rootCmd.AddCommand(projectsCmd)
}

package cmd

import (
	"fmt"
	"sort"

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
		events, err := internal.ReadAllEvents()
		if err != nil {
			return err
		}

		sessions := internal.ComputeSessions(events)

		// Extract unique project names.
		seen := make(map[string]bool)
		for _, s := range sessions {
			seen[s.Project] = true
		}
		if len(seen) == 0 {
			fmt.Println("No projects found.")
			return nil
		}

		projects := make([]string, 0, len(seen))
		for name := range seen {
			projects = append(projects, name)
		}
		sort.Strings(projects)

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

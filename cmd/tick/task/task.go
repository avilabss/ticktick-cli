package task

import (
	"github.com/avilabss/ticktick-cli/internal/ticktick"
	"github.com/spf13/cobra"
)

// NewCmd returns the "task" command group.
func NewCmd(client **ticktick.Client) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "task",
		Short: "Manage tasks",
	}
	cmd.AddCommand(listCmd(client))
	cmd.AddCommand(addCmd(client))
	cmd.AddCommand(completeCmd(client))
	cmd.AddCommand(deleteCmd(client))
	cmd.AddCommand(getCmd(client))
	return cmd
}

// NewProjectCmd returns the "project" command group.
func NewProjectCmd(client **ticktick.Client) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "project",
		Short: "Manage projects",
	}
	cmd.AddCommand(projectListCmd(client))
	return cmd
}

package task

import (
	"fmt"
	"log/slog"

	"github.com/avilabss/ticktick-cli/internal/ticktick"
	"github.com/spf13/cobra"
)

func deleteCmd(client **ticktick.Client) *cobra.Command {
	var project string

	cmd := &cobra.Command{
		Use:   "delete TASK_ID",
		Short: "Delete a task",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runDelete(*client, args[0], project)
		},
	}
	cmd.Flags().StringVar(&project, "project", "", "project ID (auto-resolved if omitted)")
	return cmd
}

func runDelete(client *ticktick.Client, taskID, project string) error {
	err := client.Task.Delete(taskID, project)
	if err != nil {
		slog.Error("Failed to delete task", "error", err)
		return err
	}

	fmt.Printf("Deleted task: %s\n", taskID)
	return nil
}

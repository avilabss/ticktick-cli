package task

import (
	"fmt"
	"log/slog"

	"github.com/avilabss/ticktick-cli/internal/ticktick"
	"github.com/spf13/cobra"
)

func completeCmd(client **ticktick.Client) *cobra.Command {
	var project string

	cmd := &cobra.Command{
		Use:   "complete TASK_ID",
		Short: "Complete a task",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runComplete(*client, args[0], project)
		},
	}
	cmd.Flags().StringVar(&project, "project", "", "project ID (auto-resolved if omitted)")
	return cmd
}

func runComplete(client *ticktick.Client, taskID, project string) error {
	_, err := client.Task.Complete(taskID, project)
	if err != nil {
		slog.Error("Failed to complete task", "error", err)
		return err
	}

	fmt.Printf("Completed task: %s\n", taskID)
	return nil
}

package task

import (
	"fmt"
	"log/slog"
	"strings"

	"github.com/avilabss/ticktick-cli/internal/ticktick"
	"github.com/spf13/cobra"
)

func getCmd(client **ticktick.Client) *cobra.Command {
	var project string

	cmd := &cobra.Command{
		Use:   "get TASK_ID",
		Short: "Get task details",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runGet(*client, args[0], project)
		},
	}
	cmd.Flags().StringVar(&project, "project", "", "project ID (auto-resolved if omitted)")
	return cmd
}

func runGet(client *ticktick.Client, taskID, project string) error {
	task, err := client.Task.Get(taskID, project)
	if err != nil {
		slog.Error("Failed to get task", "error", err)
		return err
	}

	fmt.Printf("ID:        %s\n", task.ID)
	fmt.Printf("Title:     %s\n", task.Title)
	fmt.Printf("Project:   %s\n", task.ProjectID)
	fmt.Printf("Priority:  %d\n", task.Priority)
	fmt.Printf("Status:    %s\n", formatStatus(task.Status))

	if task.Content != "" {
		fmt.Printf("Content:   %s\n", task.Content)
	}
	if len(task.Tags) > 0 {
		fmt.Printf("Tags:      %s\n", strings.Join(task.Tags, ", "))
	}
	if task.DueDate != "" {
		fmt.Printf("Due:       %s\n", task.DueDate)
	}
	if task.StartDate != "" {
		fmt.Printf("Start:     %s\n", task.StartDate)
	}
	if len(task.Items) > 0 {
		fmt.Println("Subtasks:")
		for _, item := range task.Items {
			status := "[ ]"
			if item.Status == 2 {
				status = "[x]"
			}
			fmt.Printf("  %s %s\n", status, item.Title)
		}
	}
	return nil
}

func formatStatus(status int) string {
	switch status {
	case 0:
		return "active"
	case 2:
		return "completed"
	default:
		return fmt.Sprintf("unknown (%d)", status)
	}
}

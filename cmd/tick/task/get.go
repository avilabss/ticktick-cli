package task

import (
	"flag"
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/avilabss/ticktick-cli/pkg/ticktick"
)

func runGet(client *ticktick.Client, args []string) {
	fs := flag.NewFlagSet("get", flag.ExitOnError)
	project := fs.String("project", "", "project ID (auto-resolved if omitted)")
	_ = fs.Parse(args)

	remaining := fs.Args()
	if len(remaining) == 0 {
		fmt.Println("Error: task ID is required")
		fmt.Println("Usage: tick task get TASK_ID [--project PROJECT_ID]")
		os.Exit(1)
	}

	taskID := remaining[0]
	task, err := client.Task.Get(taskID, *project)
	if err != nil {
		slog.Error("Failed to get task", "error", err)
		os.Exit(1)
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

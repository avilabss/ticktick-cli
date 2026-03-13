package task

import (
	"flag"
	"fmt"
	"log/slog"
	"os"

	"github.com/avilabss/ticktick-cli/pkg/ticktick"
)

func runComplete(client *ticktick.Client, args []string) {
	fs := flag.NewFlagSet("complete", flag.ExitOnError)
	project := fs.String("project", "", "project ID (auto-resolved if omitted)")
	_ = fs.Parse(args)

	remaining := fs.Args()
	if len(remaining) == 0 {
		fmt.Println("Error: task ID is required")
		fmt.Println("Usage: tick task complete TASK_ID [--project PROJECT_ID]")
		os.Exit(1)
	}

	taskID := remaining[0]
	_, err := client.Task.Complete(taskID, *project)
	if err != nil {
		slog.Error("Failed to complete task", "error", err)
		os.Exit(1)
	}

	fmt.Printf("Completed task: %s\n", taskID)
}

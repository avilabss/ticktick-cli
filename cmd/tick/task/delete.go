package task

import (
	"flag"
	"fmt"
	"log/slog"
	"os"

	"github.com/avilabss/ticktick-cli/pkg/ticktick"
)

func runDelete(client *ticktick.Client, args []string) {
	fs := flag.NewFlagSet("delete", flag.ExitOnError)
	project := fs.String("project", "", "project ID (auto-resolved if omitted)")
	_ = fs.Parse(args)

	remaining := fs.Args()
	if len(remaining) == 0 {
		fmt.Println("Error: task ID is required")
		fmt.Println("Usage: tick task delete TASK_ID [--project PROJECT_ID]")
		os.Exit(1)
	}

	taskID := remaining[0]
	err := client.Task.Delete(taskID, *project)
	if err != nil {
		slog.Error("Failed to delete task", "error", err)
		os.Exit(1)
	}

	fmt.Printf("Deleted task: %s\n", taskID)
}

package task

import (
	"flag"
	"fmt"
	"log/slog"
	"os"
	"strings"
	"time"

	"github.com/avilabss/ticktick-cli/pkg/ticktick"
)

func runAdd(client *ticktick.Client, args []string) {
	fs := flag.NewFlagSet("add", flag.ExitOnError)
	title := fs.String("title", "", "task title (required)")
	project := fs.String("project", "", "project name")
	tags := fs.String("tags", "", "comma-separated tags")
	priority := fs.Int("priority", 0, "priority (0=none, 1=low, 3=medium, 5=high)")
	due := fs.String("due", "", "due date (YYYY-MM-DD)")
	_ = fs.Parse(args)

	if *title == "" {
		fmt.Println("Error: --title is required")
		fmt.Println("Usage: tick task add --title \"Task name\" [--project NAME] [--tags \"a,b\"] [--priority N] [--due YYYY-MM-DD]")
		os.Exit(1)
	}

	task := ticktick.Task{
		Title:    *title,
		Priority: *priority,
		TimeZone: time.Now().Location().String(),
	}

	if *tags != "" {
		for _, t := range strings.Split(*tags, ",") {
			t = strings.TrimSpace(t)
			if t != "" {
				task.Tags = append(task.Tags, t)
			}
		}
	}

	if *due != "" {
		d, err := time.Parse("2006-01-02", *due)
		if err != nil {
			slog.Error("Invalid due date format", "error", err)
			os.Exit(1)
		}
		task.DueDate = d.Format(ticktick.TimeFormat)
		task.StartDate = task.DueDate
		task.IsAllDay = true
	}

	if *project != "" {
		projects, err := client.Task.ListProjects()
		if err != nil {
			slog.Error("Failed to list projects", "error", err)
			os.Exit(1)
		}
		found := false
		for _, p := range projects {
			if strings.Contains(strings.ToLower(p.Name), strings.ToLower(*project)) {
				task.ProjectID = p.ID
				found = true
				break
			}
		}
		if !found {
			slog.Error("Project not found", "name", *project)
			os.Exit(1)
		}
	}

	result, err := client.Task.Create(task)
	if err != nil {
		slog.Error("Failed to create task", "error", err)
		os.Exit(1)
	}

	for id := range result.ID2Etag {
		fmt.Printf("Created task: %s (ID: %s)\n", *title, id)
	}
}

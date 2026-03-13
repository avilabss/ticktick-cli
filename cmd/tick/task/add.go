package task

import (
	"fmt"
	"log/slog"
	"os"
	"strings"
	"time"

	"github.com/avilabss/ticktick-cli/internal/ticktick"
	"github.com/spf13/cobra"
)

func addCmd(client **ticktick.Client) *cobra.Command {
	var title, project, tags, due string
	var priority int

	cmd := &cobra.Command{
		Use:   "add",
		Short: "Create a new task",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runAdd(*client, title, project, tags, priority, due)
		},
	}
	cmd.Flags().StringVar(&title, "title", "", "task title (required)")
	cmd.Flags().StringVar(&project, "project", "", "project name")
	cmd.Flags().StringVar(&tags, "tags", "", "comma-separated tags")
	cmd.Flags().IntVar(&priority, "priority", 0, "priority (0=none, 1=low, 3=medium, 5=high)")
	cmd.Flags().StringVar(&due, "due", "", "due date (YYYY-MM-DD)")
	_ = cmd.MarkFlagRequired("title")
	return cmd
}

func runAdd(client *ticktick.Client, title, project, tags string, priority int, due string) error {
	task := ticktick.Task{
		Title:    title,
		Priority: priority,
		TimeZone: localTimezone(),
	}

	if tags != "" {
		for _, t := range strings.Split(tags, ",") {
			t = strings.TrimSpace(t)
			if t != "" {
				task.Tags = append(task.Tags, t)
			}
		}
	}

	if due != "" {
		d, err := time.Parse("2006-01-02", due)
		if err != nil {
			slog.Error("Invalid due date format", "error", err)
			return err
		}
		task.DueDate = d.Format(ticktick.TimeFormat)
		task.StartDate = task.DueDate
		task.IsAllDay = true
	}

	if project != "" {
		projects, err := client.Task.ListProjects()
		if err != nil {
			slog.Error("Failed to list projects", "error", err)
			return err
		}
		found := false
		for _, p := range projects {
			if strings.Contains(strings.ToLower(p.Name), strings.ToLower(project)) {
				task.ProjectID = p.ID
				found = true
				break
			}
		}
		if !found {
			slog.Error("Project not found", "name", project)
			return fmt.Errorf("project not found: %s", project)
		}
	}

	result, err := client.Task.Create(task)
	if err != nil {
		slog.Error("Failed to create task", "error", err)
		return err
	}

	for id := range result.ID2Etag {
		fmt.Printf("Created task: %s (ID: %s)\n", title, id)
	}
	return nil
}

// localTimezone returns the IANA timezone name (e.g. "Asia/Calcutta").
// Go's time.Now().Location().String() returns "Local" which the API rejects.
func localTimezone() string {
	if tz := os.Getenv("TZ"); tz != "" {
		return tz
	}
	if target, err := os.Readlink("/etc/localtime"); err == nil {
		if idx := strings.Index(target, "zoneinfo/"); idx != -1 {
			return target[idx+len("zoneinfo/"):]
		}
	}
	name, _ := time.Now().Zone()
	return name
}

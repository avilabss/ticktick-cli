package task

import (
	"fmt"
	"log/slog"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/avilabss/ticktick-cli/internal/ticktick"
	"github.com/spf13/cobra"
)

func listCmd(client **ticktick.Client) *cobra.Command {
	var project, tag string
	var priority int

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List active tasks",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runList(*client, project, tag, priority)
		},
	}
	cmd.Flags().StringVar(&project, "project", "", "filter by project name (contains, case-insensitive)")
	cmd.Flags().StringVar(&tag, "tag", "", "filter by tag (exact, case-sensitive)")
	cmd.Flags().IntVar(&priority, "priority", -1, "filter by priority (0=none, 1=low, 3=medium, 5=high)")
	return cmd
}

func runList(client *ticktick.Client, project, tag string, priority int) error {
	tasks, err := client.Task.List()
	if err != nil {
		slog.Error("Failed to list tasks", "error", err)
		return err
	}

	projects, err := client.Task.ListProjects()
	if err != nil {
		slog.Error("Failed to list projects", "error", err)
		return err
	}

	projectMap := make(map[string]string)
	for _, p := range projects {
		projectMap[p.ID] = p.Name
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	_, _ = fmt.Fprintln(w, "ID\tPROJECT\tPRI\tTITLE\tDUE\tTAGS")

	for _, t := range tasks {
		if project != "" {
			pName := strings.ToLower(projectMap[t.ProjectID])
			if !strings.Contains(pName, strings.ToLower(project)) {
				continue
			}
		}

		if tag != "" {
			found := false
			for _, tg := range t.Tags {
				if tg == tag {
					found = true
					break
				}
			}
			if !found {
				continue
			}
		}

		if priority >= 0 && t.Priority != priority {
			continue
		}

		due := "-"
		if t.DueDate != "" {
			due = t.DueDate[:10]
		}

		pName := projectMap[t.ProjectID]
		tags := strings.Join(t.Tags, ", ")

		_, _ = fmt.Fprintf(w, "%s\t%s\t%d\t%s\t%s\t%s\n",
			truncateID(t.ID), pName, t.Priority, t.Title, due, tags)
	}
	_ = w.Flush()
	return nil
}

func truncateID(id string) string {
	if len(id) > 8 {
		return id[:8]
	}
	return id
}

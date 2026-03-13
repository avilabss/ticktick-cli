package task

import (
	"flag"
	"fmt"
	"log/slog"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/avilabss/ticktick-cli/pkg/ticktick"
)

func runList(client *ticktick.Client, args []string) {
	fs := flag.NewFlagSet("list", flag.ExitOnError)
	project := fs.String("project", "", "filter by project name (contains, case-insensitive)")
	tag := fs.String("tag", "", "filter by tag (exact, case-sensitive)")
	priority := fs.Int("priority", -1, "filter by priority (0=none, 1=low, 3=medium, 5=high)")
	_ = fs.Parse(args)

	tasks, err := client.Task.List()
	if err != nil {
		slog.Error("Failed to list tasks", "error", err)
		os.Exit(1)
	}

	projects, err := client.Task.ListProjects()
	if err != nil {
		slog.Error("Failed to list projects", "error", err)
		os.Exit(1)
	}

	projectMap := make(map[string]string)
	projectIDMap := make(map[string]string)
	for _, p := range projects {
		projectMap[p.ID] = p.Name
		projectIDMap[strings.ToLower(p.Name)] = p.ID
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	_, _ = fmt.Fprintln(w, "ID\tPROJECT\tPRI\tTITLE\tDUE\tTAGS")

	for _, t := range tasks {
		if *project != "" {
			pName := strings.ToLower(projectMap[t.ProjectID])
			if !strings.Contains(pName, strings.ToLower(*project)) {
				continue
			}
		}

		if *tag != "" {
			found := false
			for _, tg := range t.Tags {
				if tg == *tag {
					found = true
					break
				}
			}
			if !found {
				continue
			}
		}

		if *priority >= 0 && t.Priority != *priority {
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
}

func truncateID(id string) string {
	if len(id) > 8 {
		return id[:8]
	}
	return id
}

package task

import (
	"fmt"
	"log/slog"
	"os"
	"text/tabwriter"

	"github.com/avilabss/ticktick-cli/pkg/ticktick"
)

func runProjectList(client *ticktick.Client) {
	projects, err := client.Task.ListProjects()
	if err != nil {
		slog.Error("Failed to list projects", "error", err)
		os.Exit(1)
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	_, _ = fmt.Fprintln(w, "ID\tNAME")

	for _, p := range projects {
		_, _ = fmt.Fprintf(w, "%s\t%s\n", p.ID, p.Name)
	}
	_ = w.Flush()
}

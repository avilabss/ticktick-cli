package task

import (
	"fmt"
	"log/slog"
	"os"
	"text/tabwriter"

	"github.com/avilabss/ticktick-cli/internal/ticktick"
	"github.com/spf13/cobra"
)

func projectListCmd(client **ticktick.Client) *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List all projects",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runProjectList(*client)
		},
	}
}

func runProjectList(client *ticktick.Client) error {
	projects, err := client.Task.ListProjects()
	if err != nil {
		slog.Error("Failed to list projects", "error", err)
		return err
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	_, _ = fmt.Fprintln(w, "ID\tNAME")

	for _, p := range projects {
		_, _ = fmt.Fprintf(w, "%s\t%s\n", p.ID, p.Name)
	}
	_ = w.Flush()
	return nil
}

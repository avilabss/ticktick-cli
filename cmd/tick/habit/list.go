package habit

import (
	"fmt"
	"log/slog"
	"os"
	"text/tabwriter"

	"github.com/avilabss/ticktick-cli/internal/ticktick"
	"github.com/spf13/cobra"
)

func listCmd(client **ticktick.Client) *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List all habits",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runList(*client)
		},
	}
}

func runList(client *ticktick.Client) error {
	habits, err := client.Habit.List()
	if err != nil {
		slog.Error("Failed to list habits", "error", err)
		return err
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	_, _ = fmt.Fprintln(w, "NAME\tTYPE\tGOAL\tCHECK-INS")

	for _, h := range habits {
		goal := fmt.Sprintf("%.0f", h.Goal)
		if h.Unit != "" && h.Unit != "Count" {
			goal = fmt.Sprintf("%.0f %s", h.Goal, h.Unit)
		}
		_, _ = fmt.Fprintf(w, "%s\t%s\t%s\t%d\n", h.Name, h.Type, goal, h.TotalCheckIns)
	}
	_ = w.Flush()
	return nil
}

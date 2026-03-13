package habit

import (
	"fmt"
	"log/slog"
	"os"
	"text/tabwriter"

	"github.com/avilabss/ticktick-cli/pkg/ticktick"
)

func runList(client *ticktick.Client, _ []string) {
	habits, err := client.Habit.List()
	if err != nil {
		slog.Error("Failed to list habits", "error", err)
		os.Exit(1)
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
}

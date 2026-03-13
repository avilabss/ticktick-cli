package habit

import (
	"fmt"
	"log/slog"
	"os"
	"text/tabwriter"
	"time"

	"github.com/avilabss/ticktick-cli/internal/ticktick"
	"github.com/spf13/cobra"
)

func statusCmd(client **ticktick.Client) *cobra.Command {
	var dateStr string

	cmd := &cobra.Command{
		Use:   "status",
		Short: "Show habit status for today",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runStatus(*client, dateStr)
		},
	}
	cmd.Flags().StringVar(&dateStr, "date", "", "date to check (YYYY-MM-DD, default today)")
	return cmd
}

func runStatus(client *ticktick.Client, dateStr string) error {
	date := time.Now()
	if dateStr != "" {
		var err error
		date, err = time.Parse("2006-01-02", dateStr)
		if err != nil {
			slog.Error("Invalid date format", "error", err)
			return err
		}
	}

	stamp := date.Year()*10000 + int(date.Month())*100 + date.Day()

	habits, err := client.Habit.List()
	if err != nil {
		slog.Error("Failed to list habits", "error", err)
		return err
	}

	if len(habits) == 0 {
		fmt.Println("No habits found")
		return nil
	}

	habitIDs := make([]string, len(habits))
	for i, h := range habits {
		habitIDs[i] = h.ID
	}

	records, err := client.Habit.GetRecords(stamp, habitIDs)
	if err != nil {
		slog.Error("Failed to get habit records", "error", err)
		return err
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	_, _ = fmt.Fprintf(w, "HABIT\tSTATUS\tVALUE/GOAL\n")

	for _, h := range habits {
		status := "Pending"
		valueStr := fmt.Sprintf("0/%.0f", h.Goal)

		if checkins, ok := records[h.ID]; ok {
			for _, c := range checkins {
				if c.CheckinStamp == stamp {
					status = "Done"
					valueStr = fmt.Sprintf("%.0f/%.0f", c.Value, h.Goal)
					break
				}
			}
		}

		_, _ = fmt.Fprintf(w, "%s\t%s\t%s\n", h.Name, status, valueStr)
	}
	_ = w.Flush()
	return nil
}

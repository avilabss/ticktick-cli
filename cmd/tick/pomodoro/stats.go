package pomodoro

import (
	"fmt"
	"log/slog"

	"github.com/avilabss/ticktick-cli/internal/ticktick"
	"github.com/spf13/cobra"
)

func statsCmd(client **ticktick.Client) *cobra.Command {
	return &cobra.Command{
		Use:   "stats",
		Short: "Show pomodoro statistics",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runStats(*client)
		},
	}
}

func runStats(client *ticktick.Client) error {
	stats, err := client.Pomodoro.Stats()
	if err != nil {
		slog.Error("Failed to get pomodoro stats", "error", err)
		return err
	}

	fmt.Printf("Today:  %d pomodoros  %s\n", stats.TodayPomoCount, formatDuration(stats.TodayPomoDuration))
	fmt.Printf("Total:  %d pomodoros  %s\n", stats.TotalPomoCount, formatDuration(stats.TotalPomoDuration))
	return nil
}

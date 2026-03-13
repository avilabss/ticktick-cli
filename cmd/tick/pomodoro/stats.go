package pomodoro

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/avilabss/ticktick-cli/pkg/ticktick"
)

func runStats(client *ticktick.Client) {
	stats, err := client.Pomodoro.Stats()
	if err != nil {
		slog.Error("Failed to get pomodoro stats", "error", err)
		os.Exit(1)
	}

	fmt.Printf("Today:  %d pomodoros  %s\n", stats.TodayPomoCount, formatDuration(stats.TodayPomoDuration))
	fmt.Printf("Total:  %d pomodoros  %s\n", stats.TotalPomoCount, formatDuration(stats.TotalPomoDuration))
}

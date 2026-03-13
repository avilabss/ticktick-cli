package pomodoro

import (
	"crypto/rand"
	"encoding/hex"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/avilabss/ticktick-cli/pkg/ticktick"
)

func runCreate(client *ticktick.Client, args []string) {
	fs := flag.NewFlagSet("create", flag.ExitOnError)
	taskID := fs.String("task", "", "task ID to associate with the pomodoro")
	duration := fs.Int("duration", 25, "pomodoro duration in minutes")
	_ = fs.Parse(args)

	now := time.Now().UTC()
	start := now.Add(-time.Duration(*duration) * time.Minute)

	b := make([]byte, 12)
	_, _ = rand.Read(b)

	pomo := ticktick.Pomodoro{
		ID:        hex.EncodeToString(b),
		StartTime: start.Format(ticktick.TimeFormat),
		EndTime:   now.Format(ticktick.TimeFormat),
		Status:    1,
		Type:      0,
	}

	if *taskID != "" {
		pomo.Tasks = []ticktick.PomodoroTask{
			{
				TaskID:    *taskID,
				StartTime: start.Format(ticktick.TimeFormat),
				EndTime:   now.Format(ticktick.TimeFormat),
			},
		}
	}

	_, err := client.Pomodoro.Create(pomo)
	if err != nil {
		slog.Error("Failed to create pomodoro", "error", err)
		os.Exit(1)
	}

	fmt.Printf("Created pomodoro: %d min\n", *duration)
}

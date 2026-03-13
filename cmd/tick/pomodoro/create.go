package pomodoro

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log/slog"
	"time"

	"github.com/avilabss/ticktick-cli/internal/ticktick"
	"github.com/spf13/cobra"
)

func createCmd(client **ticktick.Client) *cobra.Command {
	var taskID string
	var duration int

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a pomodoro record",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runCreate(*client, taskID, duration)
		},
	}
	cmd.Flags().StringVar(&taskID, "task", "", "task ID to associate with the pomodoro")
	cmd.Flags().IntVar(&duration, "duration", 25, "pomodoro duration in minutes")
	return cmd
}

func runCreate(client *ticktick.Client, taskID string, duration int) error {
	now := time.Now().UTC()
	start := now.Add(-time.Duration(duration) * time.Minute)

	b := make([]byte, 12)
	_, _ = rand.Read(b)

	pomo := ticktick.Pomodoro{
		ID:        hex.EncodeToString(b),
		StartTime: start.Format(ticktick.TimeFormat),
		EndTime:   now.Format(ticktick.TimeFormat),
		Status:    1,
		Type:      0,
	}

	if taskID != "" {
		pomo.Tasks = []ticktick.PomodoroTask{
			{
				TaskID:    taskID,
				StartTime: start.Format(ticktick.TimeFormat),
				EndTime:   now.Format(ticktick.TimeFormat),
			},
		}
	}

	_, err := client.Pomodoro.Create(pomo)
	if err != nil {
		slog.Error("Failed to create pomodoro", "error", err)
		return err
	}

	fmt.Printf("Created pomodoro: %d min\n", duration)
	return nil
}

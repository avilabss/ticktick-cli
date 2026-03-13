package habit

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/avilabss/ticktick-cli/internal/ticktick"
	"github.com/spf13/cobra"
)

func checkinCmd(client **ticktick.Client) *cobra.Command {
	var value float64

	cmd := &cobra.Command{
		Use:   "checkin HABIT_NAME",
		Short: "Check in a habit",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			habitName := strings.Join(args, " ")
			return runCheckin(*client, habitName, value)
		},
	}
	cmd.Flags().Float64Var(&value, "value", 1, "check-in value (for Number-type habits)")
	return cmd
}

func runCheckin(client *ticktick.Client, habitName string, value float64) error {
	habits, err := client.Habit.List()
	if err != nil {
		slog.Error("Failed to list habits", "error", err)
		return err
	}

	var found *ticktick.Habit
	for i, h := range habits {
		if strings.EqualFold(h.Name, habitName) {
			found = &habits[i]
			break
		}
	}
	if found == nil {
		slog.Error("Habit not found", "name", habitName)
		return fmt.Errorf("habit not found: %s", habitName)
	}

	now := time.Now()
	b := make([]byte, 12)
	_, _ = rand.Read(b)

	checkin := ticktick.HabitCheckin{
		ID:           hex.EncodeToString(b),
		HabitID:      found.ID,
		CheckinStamp: now.Year()*10000 + int(now.Month())*100 + now.Day(),
		CheckinTime:  now.UTC().Format(ticktick.TimeFormat),
		OpTime:       now.UTC().Format(ticktick.TimeFormat),
		Value:        value,
		Goal:         found.Goal,
		Status:       2,
	}

	_, err = client.Habit.Checkin(checkin)
	if err != nil {
		slog.Error("Failed to check in", "error", err)
		return err
	}

	fmt.Printf("Checked in: %s\n", found.Name)
	return nil
}

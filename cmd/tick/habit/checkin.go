package habit

import (
	"crypto/rand"
	"encoding/hex"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"strings"
	"time"

	"github.com/avilabss/ticktick-cli/pkg/ticktick"
)

func runCheckin(client *ticktick.Client, args []string) {
	fs := flag.NewFlagSet("checkin", flag.ExitOnError)
	value := fs.Float64("value", 1, "check-in value (for Number-type habits)")
	_ = fs.Parse(args)

	remaining := fs.Args()
	if len(remaining) == 0 {
		fmt.Println("Error: habit name is required")
		fmt.Println("Usage: tick habit checkin HABIT_NAME [--value N]")
		os.Exit(1)
	}

	habitName := strings.Join(remaining, " ")

	habits, err := client.Habit.List()
	if err != nil {
		slog.Error("Failed to list habits", "error", err)
		os.Exit(1)
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
		os.Exit(1)
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
		Value:        *value,
		Goal:         found.Goal,
		Status:       2,
	}

	_, err = client.Habit.Checkin(checkin)
	if err != nil {
		slog.Error("Failed to check in", "error", err)
		os.Exit(1)
	}

	fmt.Printf("Checked in: %s\n", found.Name)
}

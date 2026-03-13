package pomodoro

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/avilabss/ticktick-cli/pkg/ticktick"
)

func runDelete(client *ticktick.Client, args []string) {
	if len(args) == 0 {
		fmt.Println("Error: pomodoro ID is required")
		fmt.Println("Usage: tick pomodoro delete ID")
		os.Exit(1)
	}

	pomodoroID := args[0]

	if err := client.Pomodoro.DeletePomo(pomodoroID); err != nil {
		slog.Error("Failed to delete pomodoro", "error", err)
		os.Exit(1)
	}

	fmt.Printf("Deleted pomodoro: %s\n", pomodoroID)
}

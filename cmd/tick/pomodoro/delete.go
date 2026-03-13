package pomodoro

import (
	"fmt"
	"log/slog"

	"github.com/avilabss/ticktick-cli/internal/ticktick"
	"github.com/spf13/cobra"
)

func deleteCmd(client **ticktick.Client) *cobra.Command {
	return &cobra.Command{
		Use:   "delete ID",
		Short: "Delete a pomodoro record",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runDeletePomo(*client, args[0])
		},
	}
}

func runDeletePomo(client *ticktick.Client, pomodoroID string) error {
	if err := client.Pomodoro.DeletePomo(pomodoroID); err != nil {
		slog.Error("Failed to delete pomodoro", "error", err)
		return err
	}

	fmt.Printf("Deleted pomodoro: %s\n", pomodoroID)
	return nil
}

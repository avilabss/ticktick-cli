package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/avilabss/ticktick-cli/cmd/tick/habit"
	"github.com/avilabss/ticktick-cli/cmd/tick/pomodoro"
	"github.com/avilabss/ticktick-cli/cmd/tick/task"
	"github.com/avilabss/ticktick-cli/internal/logger"
	"github.com/avilabss/ticktick-cli/internal/ticktick"
	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
)

var (
	client  *ticktick.Client
	verbose int
)

var rootCmd = &cobra.Command{
	Use:           "tick",
	Short:         "CLI for TickTick",
	SilenceErrors: true,
	SilenceUsage:  true,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		logger.SetVerbosity(verbose)
		logger.Trace("startup", "verbosity", verbose)

		if err := godotenv.Load(); err != nil {
			slog.Error("Failed to load .env file", "error", err)
			return err
		}
		logger.Trace("Loaded .env file")

		apiToken := os.Getenv("TICKTICK_API_TOKEN")
		if apiToken == "" {
			slog.Error("TICKTICK_API_TOKEN is required")
			return fmt.Errorf("TICKTICK_API_TOKEN is required")
		}
		logger.Trace("API token loaded", "length", len(apiToken))

		var err error
		client, err = ticktick.NewTicktickClient(apiToken)
		if err != nil {
			slog.Error("Failed to create client", "error", err)
			return err
		}
		logger.Trace("Client created", "baseURL", client.BaseURL)
		return nil
	},
}

func init() {
	rootCmd.PersistentFlags().CountVarP(&verbose, "verbose", "v", "increase verbosity (-v, -vv, -vvv)")

	rootCmd.AddCommand(task.NewCmd(&client))
	rootCmd.AddCommand(task.NewProjectCmd(&client))
	rootCmd.AddCommand(pomodoro.NewCmd(&client))
	rootCmd.AddCommand(habit.NewCmd(&client))
}

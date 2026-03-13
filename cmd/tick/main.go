package main

import (
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/avilabss/ticktick-cli/cmd/tick/habit"
	"github.com/avilabss/ticktick-cli/cmd/tick/pomodoro"
	"github.com/avilabss/ticktick-cli/cmd/tick/task"
	"github.com/avilabss/ticktick-cli/pkg/logger"
	"github.com/avilabss/ticktick-cli/pkg/ticktick"
	"github.com/joho/godotenv"
)

func printUsage() {
	fmt.Println("Usage: tick [flags] <command>")
	fmt.Println()
	fmt.Println("Flags:")
	fmt.Println("  -v, -vv, -vvv    Increase verbosity level")
	fmt.Println()
	fmt.Println("Commands:")
	fmt.Println("  task        Manage tasks")
	fmt.Println("  project     List projects")
	fmt.Println("  pomodoro    Manage pomodoros")
	fmt.Println("  habit       Manage habits")
}

// parseVerbosity extracts -v/-vv/-vvv flags from args and returns
// the verbosity level and remaining args.
func parseVerbosity(args []string) (int, []string) {
	level := 0
	var remaining []string

	for _, arg := range args {
		if arg == "-v" || arg == "-vv" || arg == "-vvv" {
			l := strings.Count(arg, "v")
			if l > level {
				level = l
			}
		} else {
			remaining = append(remaining, arg)
		}
	}

	return level, remaining
}

func main() {
	verbosity, args := parseVerbosity(os.Args[1:])
	logger.SetVerbosity(verbosity)

	logger.Trace("startup", "verbosity", verbosity, "args", args)

	if err := godotenv.Load(); err != nil {
		slog.Error("Failed to load .env file", "error", err)
		os.Exit(1)
	}
	logger.Trace("Loaded .env file")

	apiToken := os.Getenv("TICKTICK_API_TOKEN")
	if apiToken == "" {
		slog.Error("TICKTICK_API_TOKEN is required")
		os.Exit(1)
	}
	logger.Trace("API token loaded", "length", len(apiToken))

	client, err := ticktick.NewTicktickClient(apiToken)
	if err != nil {
		slog.Error("Failed to create client", "error", err)
		os.Exit(1)
	}
	logger.Trace("Client created", "baseURL", client.BaseURL)

	if len(args) == 0 {
		printUsage()
		os.Exit(1)
	}

	switch args[0] {
	case "task":
		task.Run(client, args[1:])
	case "project":
		task.RunProject(client, args[1:])
	case "pomodoro":
		pomodoro.Run(client, args[1:])
	case "habit":
		habit.Run(client, args[1:])
	default:
		fmt.Printf("Unknown command: %s\n\n", args[0])
		printUsage()
		os.Exit(1)
	}
}

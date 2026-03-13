package pomodoro

import (
	"fmt"
	"os"

	"github.com/avilabss/ticktick-cli/pkg/ticktick"
)

func printUsage() {
	fmt.Println("Usage: tick pomodoro <command>")
	fmt.Println()
	fmt.Println("Commands:")
	fmt.Println("  export    Export pomodoros to CSV")
	fmt.Println("  stats     Show pomodoro statistics")
	fmt.Println("  create    Create a pomodoro record")
	fmt.Println("  delete    Delete a pomodoro record")
	fmt.Println("  timer     Manage focus timers")
}

// Run handles the "pomodoro" subcommand.
func Run(client *ticktick.Client, args []string) {
	if len(args) == 0 {
		printUsage()
		os.Exit(1)
	}

	switch args[0] {
	case "export":
		runExport(client, args[1:])
	case "stats":
		runStats(client)
	case "create":
		runCreate(client, args[1:])
	case "delete":
		runDelete(client, args[1:])
	case "timer":
		runTimer(client, args[1:])
	default:
		fmt.Printf("Unknown command: pomodoro %s\n\n", args[0])
		printUsage()
		os.Exit(1)
	}
}

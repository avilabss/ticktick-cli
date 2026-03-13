package habit

import (
	"fmt"
	"os"

	"github.com/avilabss/ticktick-cli/pkg/ticktick"
)

func printUsage() {
	fmt.Println("Usage: tick habit <command>")
	fmt.Println()
	fmt.Println("Commands:")
	fmt.Println("  list       List all habits")
	fmt.Println("  checkin    Check in a habit")
	fmt.Println("  status     Show habit status for today")
}

// Run handles the "habit" subcommand.
func Run(client *ticktick.Client, args []string) {
	if len(args) == 0 {
		printUsage()
		os.Exit(1)
	}

	switch args[0] {
	case "list":
		runList(client, args[1:])
	case "checkin":
		runCheckin(client, args[1:])
	case "status":
		runStatus(client, args[1:])
	default:
		fmt.Printf("Unknown command: habit %s\n\n", args[0])
		printUsage()
		os.Exit(1)
	}
}

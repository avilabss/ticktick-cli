package pomodoro

import (
	"fmt"
	"os"

	"github.com/avilabss/ticktick-cli/pkg/ticktick"
)

func printUsage() {
	fmt.Println("Usage: tt pomodoro <command>")
	fmt.Println()
	fmt.Println("Commands:")
	fmt.Println("  export    Export pomodoros to CSV")
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
	default:
		fmt.Printf("Unknown command: pomodoro %s\n\n", args[0])
		printUsage()
		os.Exit(1)
	}
}

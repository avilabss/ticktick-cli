package task

import (
	"fmt"
	"os"

	"github.com/avilabss/ticktick-cli/pkg/ticktick"
)

func printUsage() {
	fmt.Println("Usage: tick task <command>")
	fmt.Println()
	fmt.Println("Commands:")
	fmt.Println("  list       List active tasks")
	fmt.Println("  add        Create a new task")
	fmt.Println("  complete   Complete a task")
	fmt.Println("  delete     Delete a task")
	fmt.Println("  get        Get task details")
}

// Run handles the "task" subcommand.
func Run(client *ticktick.Client, args []string) {
	if len(args) == 0 {
		printUsage()
		os.Exit(1)
	}

	switch args[0] {
	case "list":
		runList(client, args[1:])
	case "add":
		runAdd(client, args[1:])
	case "complete":
		runComplete(client, args[1:])
	case "delete":
		runDelete(client, args[1:])
	case "get":
		runGet(client, args[1:])
	default:
		fmt.Printf("Unknown command: task %s\n\n", args[0])
		printUsage()
		os.Exit(1)
	}
}

// RunProject handles the "project" top-level command.
func RunProject(client *ticktick.Client, args []string) {
	if len(args) == 0 || args[0] == "list" {
		runProjectList(client)
	} else {
		fmt.Printf("Unknown command: project %s\n", args[0])
		fmt.Println("Usage: tick project list")
		os.Exit(1)
	}
}

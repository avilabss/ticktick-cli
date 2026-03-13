package pomodoro

import (
	"fmt"
	"log/slog"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/avilabss/ticktick-cli/pkg/ticktick"
)

func printTimerUsage() {
	fmt.Println("Usage: tick pomodoro timer <command>")
	fmt.Println()
	fmt.Println("Commands:")
	fmt.Println("  list     List all focus timers")
	fmt.Println("  stats    Show stats for a focus timer")
}

func runTimer(client *ticktick.Client, args []string) {
	if len(args) == 0 {
		printTimerUsage()
		os.Exit(1)
	}

	switch args[0] {
	case "list":
		runTimerList(client)
	case "stats":
		runTimerStats(client, args[1:])
	default:
		fmt.Printf("Unknown command: pomodoro timer %s\n\n", args[0])
		printTimerUsage()
		os.Exit(1)
	}
}

func runTimerList(client *ticktick.Client) {
	timers, err := client.Pomodoro.ListTimers()
	if err != nil {
		slog.Error("Failed to list timers", "error", err)
		os.Exit(1)
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	_, _ = fmt.Fprintln(w, "NAME\tTYPE\tDURATION")

	for _, t := range timers {
		_, _ = fmt.Fprintf(w, "%s\t%s\t%dm\n", t.Name, t.Type, t.PomodoroTime)
	}
	_ = w.Flush()
}

func runTimerStats(client *ticktick.Client, args []string) {
	if len(args) == 0 {
		fmt.Println("Error: timer name is required")
		fmt.Println("Usage: tick pomodoro timer stats NAME")
		os.Exit(1)
	}

	timerName := strings.Join(args, " ")

	timers, err := client.Pomodoro.ListTimers()
	if err != nil {
		slog.Error("Failed to list timers", "error", err)
		os.Exit(1)
	}

	var found *ticktick.FocusTimer
	for i, t := range timers {
		if strings.EqualFold(t.Name, timerName) {
			found = &timers[i]
			break
		}
	}
	if found == nil {
		slog.Error("Timer not found", "name", timerName)
		os.Exit(1)
	}

	overview, err := client.Pomodoro.TimerOverview(found.ID)
	if err != nil {
		slog.Error("Failed to get timer stats", "error", err)
		os.Exit(1)
	}

	fmt.Printf("Timer:  %s\n", found.Name)
	fmt.Printf("Days:   %d\n", overview.Days)
	fmt.Printf("Today:  %s\n", formatDuration(overview.Today))
	fmt.Printf("Total:  %s\n", formatDuration(overview.Total))
}

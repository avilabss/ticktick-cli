package pomodoro

import (
	"fmt"
	"log/slog"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/avilabss/ticktick-cli/internal/ticktick"
	"github.com/spf13/cobra"
)

func timerCmd(client **ticktick.Client) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "timer",
		Short: "Manage focus timers",
	}
	cmd.AddCommand(timerListCmd(client))
	cmd.AddCommand(timerStatsCmd(client))
	return cmd
}

func timerListCmd(client **ticktick.Client) *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List all focus timers",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runTimerList(*client)
		},
	}
}

func timerStatsCmd(client **ticktick.Client) *cobra.Command {
	return &cobra.Command{
		Use:   "stats NAME",
		Short: "Show stats for a focus timer",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			timerName := strings.Join(args, " ")
			return runTimerStats(*client, timerName)
		},
	}
}

func runTimerList(client *ticktick.Client) error {
	timers, err := client.Pomodoro.ListTimers()
	if err != nil {
		slog.Error("Failed to list timers", "error", err)
		return err
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	_, _ = fmt.Fprintln(w, "NAME\tTYPE\tDURATION")

	for _, t := range timers {
		_, _ = fmt.Fprintf(w, "%s\t%s\t%dm\n", t.Name, t.Type, t.PomodoroTime)
	}
	_ = w.Flush()
	return nil
}

func runTimerStats(client *ticktick.Client, timerName string) error {
	timers, err := client.Pomodoro.ListTimers()
	if err != nil {
		slog.Error("Failed to list timers", "error", err)
		return err
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
		return fmt.Errorf("timer not found: %s", timerName)
	}

	overview, err := client.Pomodoro.TimerOverview(found.ID)
	if err != nil {
		slog.Error("Failed to get timer stats", "error", err)
		return err
	}

	fmt.Printf("Timer:  %s\n", found.Name)
	fmt.Printf("Days:   %d\n", overview.Days)
	fmt.Printf("Today:  %s\n", formatDuration(overview.Today))
	fmt.Printf("Total:  %s\n", formatDuration(overview.Total))
	return nil
}

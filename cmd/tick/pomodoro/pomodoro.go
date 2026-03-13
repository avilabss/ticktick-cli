package pomodoro

import (
	"github.com/avilabss/ticktick-cli/internal/ticktick"
	"github.com/spf13/cobra"
)

// NewCmd returns the "pomodoro" command group.
func NewCmd(client **ticktick.Client) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "pomodoro",
		Short: "Manage pomodoros",
	}
	cmd.AddCommand(exportCmd(client))
	cmd.AddCommand(statsCmd(client))
	cmd.AddCommand(createCmd(client))
	cmd.AddCommand(deleteCmd(client))
	cmd.AddCommand(timerCmd(client))
	return cmd
}

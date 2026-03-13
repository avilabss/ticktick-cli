package habit

import (
	"github.com/avilabss/ticktick-cli/internal/ticktick"
	"github.com/spf13/cobra"
)

// NewCmd returns the "habit" command group.
func NewCmd(client **ticktick.Client) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "habit",
		Short: "Manage habits",
	}
	cmd.AddCommand(listCmd(client))
	cmd.AddCommand(checkinCmd(client))
	cmd.AddCommand(statusCmd(client))
	return cmd
}

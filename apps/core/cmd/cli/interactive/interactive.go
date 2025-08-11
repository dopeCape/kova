package interactive

import (
	"github.com/dopeCape/kova/internal/tui"
	"github.com/spf13/cobra"
)

var (
	InteractiveCmd = &cobra.Command{
		Use:   "interactive",
		Short: "Runs kova-cli in interactive mode",
		RunE: func(cmd *cobra.Command, args []string) error {
			return tui.StartTUI()
		},
	}
)

func Execute() error {
	return InteractiveCmd.Execute()
}

func init() {
}

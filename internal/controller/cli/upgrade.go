package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

func (c CliController) newUpgradeCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "upgrade",
		Short: "Upgrade Flow to latest version",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("upgrade")

			return nil
		},
	}

	return cmd
}

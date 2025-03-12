package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

func (c CliController) newInitCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "init",
		Short: "Creates a new config file in project directory",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("init")

			return nil
		},
	}

	return cmd
}

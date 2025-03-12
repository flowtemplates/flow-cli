package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

var Version = "dev"

func (c CliController) newVersionCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "version",
		Short: "Print Flow version",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println(Version)

			return nil
		},
	}

	return cmd
}

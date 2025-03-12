package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

func (c CliController) newCloneCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "clone",
		Short: "Creates template with <template_name> from directory or file located in path",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("clone")

			return nil
		},
	}

	return cmd
}

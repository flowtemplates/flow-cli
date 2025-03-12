package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

func (c CliController) newLspProxyCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "lsp-proxy",
		Short: "Start a server for the Language Server Protocol over stdin/stdout",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("lsp")

			return nil
		},
	}

	return cmd
}

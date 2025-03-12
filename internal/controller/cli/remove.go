package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

func (c CliController) newRemoveCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "remove <template name>",
		Short:   "Remove template by name",
		Args:    cobra.MinimumNArgs(1),
		Aliases: []string{"rm"},
		RunE: func(cmd *cobra.Command, args []string) error {
			templateName := args[0]

			// err := c.service.Remove(templateName)
			// if err != nil {
			// 	return fmt.Errorf("failed to add: %w", err)
			// }
			fmt.Println(templateName)

			return nil
		},
	}

	cmd.Flags().Bool("print-json", false, "Output in JSON format")

	return cmd
}

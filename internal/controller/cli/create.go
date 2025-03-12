package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

func (c CliController) newCreateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create <template name> [...paths]",
		Short: "Create selected template to output dirs",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			templateName := args[0]
			vars := c.parseVars(args[1:])
			paths, _ := cmd.Flags().GetStringSlice("values")

			overWriteFn := func(path []string) ([]string, error) {
				return []string{}, nil
			}

			err := c.service.Create(templateName, vars, overWriteFn, paths...)
			if err != nil {
				return fmt.Errorf("failed to add: %w", err)
			}

			return nil
		},
	}

	cmd.Flags().StringSliceP("values", "v", []string{}, "Values to pass to context")
	cmd.Flags().Bool("print-json", false, "Output in JSON format")

	return cmd
}

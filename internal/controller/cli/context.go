package cli

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
)

func (c CliController) newContextCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "context",
		Short: "Print all variables with corresponding types available for template (context)",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			templateName := args[0]
			tm, err := c.service.GetTemplateContext(templateName)
			if err != nil {
				return fmt.Errorf("failed to get template: %w", err)
			}
			j, _ := json.Marshal(tm)
			fmt.Printf("%s\n", j)

			return nil
		},
	}

	// cmd.Flags().Bool("json", false, "Output in JSON format")
	return cmd
}

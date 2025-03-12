package cli

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
)

func (c CliController) newListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all template names",
		RunE: func(cmd *cobra.Command, args []string) error {
			jsonOutput, _ := cmd.Flags().GetBool("json")

			templates, err := c.service.ListTemplates()
			if err != nil {
				return fmt.Errorf("failed to list templates: %w", err)
			}

			if jsonOutput {
				data, _ := json.Marshal(templates)
				fmt.Printf("%s\n", data)
			} else {
				for _, c := range templates {
					fmt.Printf("- %s\n", c)
				}
			}
			return nil
		},
	}

	cmd.Flags().Bool("print-json", false, "Output in JSON format")
	return cmd
}

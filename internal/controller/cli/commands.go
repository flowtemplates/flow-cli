package cli

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
)

func (c CliController) newListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all components",
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

	cmd.Flags().Bool("json", false, "Output in JSON format")
	return cmd
}

func (c CliController) newAddCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add",
		Short: "Add component to dir",
		Args:  cobra.MinimumNArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			templateName := args[0]
			err := c.service.Add(templateName, args[1:]...)
			if err != nil {
				return fmt.Errorf("failed to add: %w", err)
			}

			return nil
		},
	}

	cmd.Flags().Bool("json", false, "Output in JSON format")
	return cmd
}

func (c CliController) newGetCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get",
		Short: "Get component",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			// fmt.Println(args)
			templateName := args[0]
			tm, err := c.service.Get(templateName)
			if err != nil {
				return fmt.Errorf("failed to get template: %w", err)
			}
			j, _ := json.Marshal(tm)
			fmt.Printf("%s\n", j)

			return nil
		},
	}

	cmd.Flags().Bool("json", false, "Output in JSON format")
	return cmd
}

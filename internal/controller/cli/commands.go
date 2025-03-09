package cli

import (
	"encoding/json"
	"fmt"
	"strings"

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

func (c CliController) parseVars(vars []string) map[string]*string {
	res := make(map[string]*string)
	for _, v := range vars {
		if strings.Contains(v, "=") {
			parts := strings.SplitN(v, "=", 2)
			res[parts[0]] = &parts[1]
		} else {
			res[v] = nil
		}
	}

	return res
}

func (c CliController) newAddCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add <template name> [flags] [destination paths...]",
		Short: "Add component to dirs",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			templateName := args[0]
			vars := c.parseVars(args[1:])
			// fmt.Println("templateName:", templateName)
			// fmt.Println("vars:", vars)

			paths, _ := cmd.Flags().GetStringSlice("out")
			// fmt.Println("Paths:", paths)

			overWriteFn := func(path string) bool {
				fmt.Printf("overwrite %s\n", path)
				return false
			}

			err := c.service.Add(templateName, vars, overWriteFn, paths...)
			if err != nil {
				return fmt.Errorf("failed to add: %w", err)
			}

			return nil
		},
	}

	cmd.Flags().StringSliceP("out", "o", []string{}, "Output paths")
	cmd.Flags().Bool("json", false, "Output in JSON format")

	return cmd
}

func (c CliController) newGetCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get",
		Short: "Get component",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
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

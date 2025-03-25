package main

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"slices"
	"strings"

	"github.com/flowtemplates/flow-cli/internal/config"
	"github.com/flowtemplates/flow-cli/internal/lsp"
	"github.com/flowtemplates/flow-cli/internal/repository/source"
	"github.com/flowtemplates/flow-cli/internal/repository/templates"
	"github.com/flowtemplates/flow-cli/internal/service"
	"github.com/spf13/cobra"
)

func createService() (*service.Service, error) {
	cfg, err := config.GetConfig(defaultConfigName)
	if err != nil {
		return nil, fmt.Errorf("failed to get config: %w", err)
	}

	tr := templates.New(cfg.TemplatesFolder)
	sr := source.New()

	return service.New(tr, sr), nil
}

type configExt string

const (
	configJson  configExt = "json"
	configJsonc configExt = "jsonc"
	configYaml  configExt = "yaml"
	configYml   configExt = "yml"
)

var configExts = []string{
	string(configJson),
	string(configJsonc),
	string(configYaml),
	string(configYml),
}

func (e *configExt) String() string {
	return string(*e)
}

func (e *configExt) Set(v string) error {
	if slices.Contains(configExts, v) {
		*e = configExt(v)
		return nil
	}

	return fmt.Errorf("must be one of: %s", strings.Join(configExts, ", "))
}

func (e *configExt) Type() string {
	return "string"
}

func newInitCmd() *cobra.Command {
	ext := configYml
	cmd := &cobra.Command{
		Use:   "init",
		Short: "Creates a new config file in project directory",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println(ext.String())
			return nil
		},
	}

	cmd.Flags().BoolP("path", "p", false, "Path to directory")
	cmd.Flags().VarP(&ext, "ext", "e", "Config file extension")

	return cmd
}

func newListCmd() *cobra.Command {
	var printJson bool
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all template names",
		RunE: func(cmd *cobra.Command, args []string) error {
			s, err := createService()
			if err != nil {
				return err
			}

			templates, err := s.ListTemplates()
			if err != nil {
				return fmt.Errorf("failed to list templates: %w", err)
			}

			if printJson {
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

	cmd.Flags().BoolVar(&printJson, "print-json", false, "Output in JSON format")
	return cmd
}

func newCloneCmd() *cobra.Command {
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

func newLspProxyCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "lsp-proxy",
		Short: "Start a server for the Language Server Protocol over stdin/stdout",
		RunE: func(cmd *cobra.Command, args []string) error {
			logFilePath := "/home/skewbik/dev/seriousbiz/flowtemplates/flow-cli/.out/log"
			file, err := os.OpenFile(logFilePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o644)
			if err != nil {
				return fmt.Errorf("failed to open log-file: %w", err)
			}
			defer file.Close()

			logger := slog.New(slog.NewJSONHandler(file, &slog.HandlerOptions{
				Level: slog.LevelDebug,
			}))

			server := lsp.NewServer(&lsp.ServerOptions{
				In:     os.Stdin,
				Out:    os.Stdout,
				Err:    os.Stderr,
				Logger: logger,
			})

			return server.Run(cmd.Context())
		},
	}

	return cmd
}

func newCreateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "create <template name> [...paths]",
		Short:   "Create selected template to output dirs",
		Args:    cobra.MinimumNArgs(1),
		Aliases: []string{"c"},
		RunE: func(cmd *cobra.Command, args []string) error {
			templateName := args[0]
			paths := args[1:]
			values, _ := cmd.Flags().GetStringSlice("values")
			vars := parseVars(values)

			overWriteFn := func(p []string) ([]string, error) {
				fmt.Printf("request to overwrite: %v\n", p)
				return []string{}, nil
			}

			s, err := createService()
			if err != nil {
				return err
			}

			if err := s.Create(templateName, vars, overWriteFn, paths...); err != nil {
				return fmt.Errorf("failed to add: %w", err)
			}

			return nil
		},
	}

	cmd.Flags().StringSliceP("values", "v", []string{}, "Values to pass to context")
	cmd.Flags().Bool("print-json", false, "Output in JSON format")

	return cmd
}

var Version = "dev"

func newVersionCmd() *cobra.Command {
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

func newRemoveCmd() *cobra.Command {
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

func newUpgradeCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "upgrade",
		Short: "Upgrade Flow to latest version",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("upgrade")

			return nil
		},
	}

	return cmd
}

func newContextCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "context",
		Short: "Print all variables with corresponding types available for template (context)",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			templateName := args[0]
			s, err := createService()
			if err != nil {
				return err
			}

			tm, err := s.GetTemplateContext(templateName)
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

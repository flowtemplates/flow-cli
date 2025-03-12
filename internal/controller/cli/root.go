package cli

import (
	"fmt"

	"github.com/charmbracelet/huh"
	"github.com/flowtemplates/flow-go/types"
	"github.com/spf13/cobra"
)

func (c CliController) Cmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "flow",
		Short: "Flow CLI",
		Long:  "Modern toolchain for component code generation.",
		RunE: func(_ *cobra.Command, _ []string) error {
			return c.handleMain()
		},
	}

	rootCmd.AddCommand(c.newListCmd())
	rootCmd.AddCommand(c.newContextCmd())
	rootCmd.AddCommand(c.newCreateCmd())
	rootCmd.AddCommand(c.newRemoveCmd())
	rootCmd.AddCommand(c.newUpgradeCmd())
	rootCmd.AddCommand(c.newInitCmd())
	rootCmd.AddCommand(c.newCloneCmd())
	rootCmd.AddCommand(c.newVersionCmd())
	rootCmd.AddCommand(c.newLspProxyCmd())

	return rootCmd
}

func (c CliController) handleMain() error {
	templates, err := c.service.ListTemplates()
	if err != nil {
		return fmt.Errorf("failed to load templates: %w", err)
	}

	var templateName string

	templateForm := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Select a template").
				OptionsFunc(func() []huh.Option[string] {
					var options []huh.Option[string]
					for _, t := range templates {
						options = append(options, huh.NewOption(t, t))
					}
					return options
				}, &templateName).
				Value(&templateName),
		),
	)

	if err := templateForm.Run(); err != nil {
		return fmt.Errorf("failed to run template form: %w", err)
	}

	tm, err := c.service.GetTemplateContext(templateName)
	if err != nil {
		return fmt.Errorf("failed to get template: %w", err)
	}

	variableMap := make(map[string]*string)

	var formFields []huh.Field
	var flagFields []huh.Option[string]

	for name, typ := range tm {
		if typ == types.Boolean {
			flagFields = append(flagFields, huh.NewOption(name, name))
		} else {
			var input string
			formFields = append(formFields, huh.NewInput().
				Title(name).
				Key(name).
				Value(&input),
			)

			variableMap[name] = &input
		}
	}

	var selectedFlags []string
	var dest string

	groups := []*huh.Group{}
	if len(formFields) > 0 {
		groups = append(groups, huh.NewGroup(formFields...))
	}

	if len(flagFields) > 0 {
		groups = append(groups, huh.NewGroup(
			huh.NewMultiSelect[string]().
				Options(flagFields...).
				Title("Select flags").
				Value(&selectedFlags),
		))
	}

	groups = append(groups,
		huh.NewGroup(
			huh.NewFilePicker().
				DirAllowed(true).
				FileAllowed(false).
				Height(10).
				Picking(true).
				ShowPermissions(false).
				Value(&dest),
		))

	paramsForm := huh.NewForm(groups...)

	if err := paramsForm.Run(); err != nil {
		return fmt.Errorf("failed to run form: %w", err)
	}

	for _, name := range selectedFlags {
		variableMap[name] = nil
	}

	overWriteFn := func(paths []string) ([]string, error) {
		ov := []string{}
		overwriteForm := huh.NewForm(
			huh.NewGroup(
				huh.NewMultiSelect[string]().
					Title("Select files to overwrite").
					OptionsFunc(func() []huh.Option[string] {
						var options []huh.Option[string]
						for _, t := range paths {
							options = append(options, huh.NewOption(t, t).Selected(true))
						}
						return options
					}, &templateName).
					Value(&ov),
			),
		)
		if err := overwriteForm.Run(); err != nil {
			return nil, fmt.Errorf("failed to run overwrite form: %w", err)
		}

		return ov, nil
	}

	if err := c.service.Create(templateName, variableMap, overWriteFn, dest); err != nil {
		return fmt.Errorf("failed to add: %w", err)
	}

	return nil
}

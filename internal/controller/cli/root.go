package cli

import (
	"fmt"
	"log/slog"

	"github.com/charmbracelet/huh"
	"github.com/flowtemplates/cli/pkg/flow-go/analyzer"
	"github.com/flowtemplates/cli/pkg/flow-go/types"
	"github.com/spf13/cobra"
)

type iService interface {
	ListTemplates() ([]string, error)
	Add(templateName string, scope map[string]*string, overwriteFn func(path string) bool, dests ...string) error
	Get(templateName string) (analyzer.TypeMap, error)
}

type CliController struct {
	service iService
	logger  *slog.Logger
}

func New(service iService, logger *slog.Logger) *CliController {
	return &CliController{
		service: service,
		logger:  logger,
	}
}

func (c CliController) Cmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "flow",
		Short: "FlowTemplates CLI",
		RunE: func(cmd *cobra.Command, args []string) error {
			return c.showForm()
		},
	}

	rootCmd.AddCommand(c.newListCmd())
	rootCmd.AddCommand(c.newGetCmd())
	rootCmd.AddCommand(c.newAddCmd())

	return rootCmd
}

func (c CliController) showForm() error {
	var templateName string

	templates, err := c.service.ListTemplates()
	if err != nil {
		return fmt.Errorf("failed to load templates: %w", err)
	}

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
		return err
	}

	tm, err := c.service.Get(templateName)
	if err != nil {
		return fmt.Errorf("failed to get template: %w", err)
	}

	variableMap := make(map[string]*string)

	var formFields []huh.Field
	var flagFields []huh.Option[bool]

	for name, typ := range tm {
		if typ == types.Boolean {
			flagFields = append(flagFields, huh.NewOption(name, true))
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

	var res []bool
	var dest string

	groups := []*huh.Group{}
	if len(formFields) > 0 {
		groups = append(groups, huh.NewGroup(formFields...))
	}

	if len(flagFields) > 0 {
		groups = append(groups, huh.NewGroup(
			huh.NewMultiSelect[bool]().
				Options(flagFields...).
				Title("Flags").
				Value(&res),
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

	dynamicForm := huh.NewForm(groups...)

	if err := dynamicForm.Run(); err != nil {
		return err
	}

	// fmt.Printf("vars: %v\n", variableMap)
	// fmt.Printf("res: %v\n", res)
	// fmt.Printf("dest: %v\n", dest)

	overWriteFn := func(path string) bool {
		// fmt.Printf("overwrite %s ", path)
		return false
	}

	if err := c.service.Add(templateName, variableMap, overWriteFn, dest); err != nil {
		return fmt.Errorf("failed to add: %w", err)
	}

	return nil
}

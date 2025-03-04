package cli

import (
	"log/slog"

	"github.com/spf13/cobra"
)

type iService interface {
	ListTemplates() ([]string, error)
	Add(templateName string, dests ...string) error
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
	}

	rootCmd.AddCommand(c.newListCmd())
	rootCmd.AddCommand(c.newGetCmd())
	rootCmd.AddCommand(c.newAddCmd())

	return rootCmd
}

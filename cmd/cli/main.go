package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"os"

	"github.com/flowtemplates/cli/internal/config"
	"github.com/flowtemplates/cli/internal/controller/cli"
	"github.com/flowtemplates/cli/internal/repository/templates"
	"github.com/flowtemplates/cli/internal/service"
)

const defaultConfigName = "flow"

func run() error {
	ctx := context.Background()

	cfg, err := config.GetConfig(defaultConfigName)
	if err != nil {
		return fmt.Errorf("failed to get config: %w", err)
	}

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	r := templates.New(cfg.TemplatesFolder)
	s := service.New(r)
	c := cli.New(s, logger)

	return c.Cmd().ExecuteContext(ctx) // nolint: wrapcheck
}

func main() {
	if err := run(); err != nil {
		log.Fatalln(err)
	}
}

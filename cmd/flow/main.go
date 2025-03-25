package main

import (
	"context"
	"fmt"
	"os"
)

const defaultConfigName = "flow"

// func run() error {

// 	cfg, err := config.GetConfig(defaultConfigName)
// 	if err != nil {
// 		return fmt.Errorf("failed to get config: %w", err)
// 	}

// 	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

// 	tr := templates.New(cfg.TemplatesFolder)
// 	sr := source.New()

// 	s := service.New(tr, sr)
// 	c := cli.New(s, logger)

// 	return
// }

func main() {
	ctx := context.Background()
	if err := cmd().ExecuteContext(ctx); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

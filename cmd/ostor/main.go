package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/Luzilla/acronis-s3-usage/internal/cmd"
	"github.com/urfave/cli/v3"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	// initialize a default logger
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelDebug})))

	app := &cli.Command{
		Name: "ostor-client",
		Authors: []any{
			map[string]string{
				"Name": "Luzilla Capital GmbH",
			},
		},
		Usage:   "a program to interact with the s3 management APIs in ACI and VHI",
		Version: fmt.Sprintf("%s (%s, date: %s)", version, commit, date),
		Before:  cmd.Before,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "s3-endpoint",
				Sources:  cli.EnvVars("S3_ENDPOINT"),
				Required: true,
			},
			&cli.StringFlag{
				Name:     "s3-system-key-id",
				Sources:  cli.EnvVars("S3_SYSTEM_KEY_ID"),
				Required: true,
			},
			&cli.StringFlag{
				Name:     "s3-system-secret",
				Sources:  cli.EnvVars("S3_SYSTEM_SECRET_KEY"),
				Required: true,
			},
		},
		Commands: []*cli.Command{
			cmd.BucketCommand(),
			cmd.StatsCommand(),
			cmd.UsersCommand(),
		},
	}

	if err := app.Run(context.Background(), os.Args); err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
}

package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/Luzilla/acronis-s3-usage/internal/cmd"
	"github.com/Luzilla/acronis-s3-usage/pkg/ostor"
	"github.com/urfave/cli/v2"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	// initialize a default logger
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelDebug})))

	app := &cli.App{
		Name:     "ostor-client",
		HelpName: "a program to interact with the s3 management APIs in ACI and VHI",
		Version:  fmt.Sprintf("%s (%s, date: %s)", version, commit, date),
		Before: func(cCtx *cli.Context) error {
			client, err := ostor.New(
				cCtx.String("s3-endpoint"),
				cCtx.String("s3-system-key-id"),
				cCtx.String("s3-system-secret"))
			if err != nil {
				return err
			}

			cCtx.Context = context.WithValue(cCtx.Context, cmd.OstorClient, client)
			return nil
		},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "s3-endpoint",
				EnvVars:  []string{"S3_ENDPOINT"},
				Required: true,
			},
			&cli.StringFlag{
				Name:     "s3-system-key-id",
				EnvVars:  []string{"S3_SYSTEM_KEY_ID"},
				Required: true,
			},
			&cli.StringFlag{
				Name:     "s3-system-secret",
				EnvVars:  []string{"S3_SYSTEM_SECRET_KEY"},
				Required: true,
			},
		},
		Commands: []*cli.Command{
			{
				Name:    "buckets",
				Aliases: []string{"b"},
				Usage:   "list buckets",
				Action:  cmd.ListBuckets,
				Flags: []cli.Flag{
					emailFlag(),
				},
			},
			{
				Name:    "stats",
				Aliases: []string{"s"},
				Usage:   "list stats",
				Flags:   []cli.Flag{
					// &cli.TimestampFlag{
					// 	Name:     "from",
					// 	Layout:   "2006-01-02",
					// 	Timezone: time.UTC,
					// 	Required: false,
					// },
				},
				Action: cmd.ShowStats,
			},
			{
				Name:    "users",
				Aliases: []string{"u"},
				Usage:   "manage users",
				Action:  cmd.Users,
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:  "usage",
						Value: false,
					},
				},
				Subcommands: []*cli.Command{
					{
						Name: "create",
						Flags: []cli.Flag{
							emailFlag(),
						},
						Action: cmd.CreateUser,
					},
					{
						Name: "lock",
						Flags: []cli.Flag{
							emailFlag(),
						},
						Action: cmd.LockUser,
					},
					{
						Name: "unlock",
						Flags: []cli.Flag{
							emailFlag(),
						},
						Action: cmd.UnlockUser,
					},
					{
						Name: "show",
						Flags: []cli.Flag{
							emailFlag(),
						},
						Action: cmd.ShowUser,
					},
					{
						Name: "create-key",
						Flags: []cli.Flag{
							emailFlag(),
						},
						Action: cmd.CreateKey,
					},
					{
						Name: "revoke-key",
						Flags: []cli.Flag{
							emailFlag(),
							&cli.StringFlag{
								Name:     "key-id",
								Required: true,
							},
						},
						Action: cmd.RevokeKey,
					},
					{
						Name: "limits",
						Flags: []cli.Flag{
							emailFlag(),
						},
						Action: cmd.UserLimits,
					},
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
}

func emailFlag() *cli.StringFlag {
	return &cli.StringFlag{
		Name:     "email",
		Required: true,
	}
}

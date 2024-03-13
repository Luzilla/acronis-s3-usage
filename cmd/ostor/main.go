package main

import (
	"context"
	"fmt"
	"log"
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
			},
			{
				Name:    "stats",
				Aliases: []string{"s"},
				Usage:   "list stats",
				Action:  cmd.List,
			},
			{
				Name:    "users",
				Aliases: []string{"u"},
				Usage:   "list users",
				Action:  cmd.Users,
				Subcommands: []*cli.Command{
					{
						Name: "show",
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:     "email",
								Required: true,
							},
						},
						Action: cmd.ShowUser,
					},
					{
						Name: "create-key",
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:     "email",
								Required: true,
							},
						},
						Action: cmd.CreateKey,
					},
					{
						Name: "revoke-key",
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:     "email",
								Required: true,
							},
							&cli.StringFlag{
								Name:     "key-id",
								Required: true,
							},
						},
						Action: cmd.RevokeKey,
					},
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

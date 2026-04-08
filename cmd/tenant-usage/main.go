package main

import (
	"context"
	"fmt"
	"log"
	"math"
	"os"

	"github.com/Luzilla/acronis-s3-usage/internal/utils"
	"github.com/Luzilla/acronis-s3-usage/pkg/acronis"
	"github.com/urfave/cli/v3"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	app := &cli.Command{
		Name:    "tenant-usgae",
		Usage:   "a program to interact with the ACI APIs to extract s3 basic usage",
		Version: fmt.Sprintf("%s (%s, date: %s)", version, commit, date),
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "client-id",
				Required: true,
				Sources:  cli.EnvVars("ACI_CLIENT_ID"),
			},
			&cli.StringFlag{
				Name:     "secret",
				Required: true,
				Sources:  cli.EnvVars("ACI_SECRET"),
			},
			&cli.StringFlag{
				Name:     "dc-url",
				Required: true,
				Sources:  cli.EnvVars("ACI_DC_URL"),
				Value:    "https://eu2-cloud.acronis.com",
			},
		},
		Action: func(ctx context.Context, c *cli.Command) error {
			aci := acronis.NewClient(
				c.String("client-id"),
				c.String("secret"),
				c.String("dc-url"),
			)

			tenantId, err := aci.GetTenantID(ctx)
			if err != nil {
				return err
			}
			fmt.Printf("Got tenant id: %s\n\n", tenantId)

			usageData, err := aci.GetUsage(ctx, tenantId)
			if err != nil {
				return err
			}

			for _, items := range usageData.Items {
				for _, usages := range items.Usages {
					if usages.Name != "hci_s3_storage" {
						continue
					}

					app, err := aci.GetApplication(ctx, usages.ApplicationID)
					if err != nil {
						panic(err)
					}

					fmt.Printf("%s (Type: %s)\n%s -- %s\n\n",
						app.Name,
						app.Type,
						usages.Name,
						utils.PrettyByteSize(int64(math.Round(usages.AbsoluteValue))),
					)
				}
			}

			return nil
		},
	}

	if err := app.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}

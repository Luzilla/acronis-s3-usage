package main

import (
	"fmt"
	"log"
	"math"
	"os"

	"github.com/Luzilla/acronis-s3-usage/internal/utils"
	"github.com/Luzilla/acronis-s3-usage/pkg/acronis"
	"github.com/urfave/cli/v2"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	app := &cli.App{
		Name:     "tenant-usgae",
		HelpName: "a program to interact with the ACI APIs to extract s3 basic usage",
		Version:  fmt.Sprintf("%s (%s, date: %s)", version, commit, date),
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "client-id",
				Required: true,
				EnvVars:  []string{"ACI_CLIENT_ID"},
			},
			&cli.StringFlag{
				Name:     "secret",
				Required: true,
				EnvVars:  []string{"ACI_SECRET"},
			},
			&cli.StringFlag{
				Name:     "dc-url",
				Required: true,
				EnvVars:  []string{"ACI_DC_URL"},
				Value:    "https://eu2-cloud.acronis.com",
			},
		},
		Action: func(cCtx *cli.Context) error {
			aci := acronis.NewClient(
				cCtx.String("client-id"),
				cCtx.String("secret"),
				cCtx.String("dc-url"),
			)

			tenantId, err := aci.GetTenantID()
			if err != nil {
				return err
			}
			fmt.Printf("Got tenant id: %s\n\n", tenantId)

			usageData, err := aci.GetUsage(tenantId)
			if err != nil {
				return err
			}

			for _, items := range usageData.Items {
				for _, usages := range items.Usages {
					if usages.Name != "hci_s3_storage" {
						continue
					}

					app, err := aci.GetApplication(usages.ApplicationID)
					if err != nil {
						panic(err)
					}

					fmt.Printf("%s (Type: %s)\n%s -- %s\n\n",
						app.Name,
						app.Type,
						usages.Name,
						utils.PrettyByteSize(int(math.Round(usages.AbsoluteValue))),
					)
				}
			}

			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

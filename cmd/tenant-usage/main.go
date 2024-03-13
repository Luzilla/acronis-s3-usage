package main

import (
	"fmt"
	"math"
	"os"

	"github.com/Luzilla/acronis-s3-usage/internal/utils"
	"github.com/Luzilla/acronis-s3-usage/pkg/acronis"
)

func main() {
	aci := acronis.NewClient(
		os.Getenv("ACI_CLIENT_ID"),
		os.Getenv("ACI_SECRET"),
		os.Getenv("ACI_DC_URL"),
	)

	tenantId, err := aci.GetTenantID()
	if err != nil {
		panic(err)
	}
	fmt.Printf("Got tenant id: %s\n\n", tenantId)

	usageData, err := aci.GetUsage(tenantId)
	if err != nil {
		panic(err)
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
}

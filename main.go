package main

import (
	"fmt"
	"os"

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
	fmt.Printf("Got tenant id: %s\n", tenantId)

	usageData, err := aci.GetUsage(tenantId)
	if err != nil {
		panic(err)
	}

	for _, items := range usageData.Items {
		//fmt.Printf("Got tenant ID: %s\n", items.Tenant)
		for _, usages := range items.Usages {
			if usages["name"] != "hci_s3_storage" {
				continue
			}

			app, err := aci.GetApplication(usages["application_id"].(string))
			if err != nil {
				panic(err)
			}
			fmt.Printf("%s (Type: %s)\n\n%s -- %.2f GB\n",
				app.Name,
				app.Type,
				usages["name"],
				// bitshift -> byte to gb
				(usages["absolute_value"].(float64) / (1 << 30)))
		}
	}
}

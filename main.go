package main

import (
	"encoding/base64"
	"fmt"
	"os"

	"github.com/go-resty/resty/v2"
)

func encodeClientCredentials() string {
	authStr := fmt.Sprintf("%s:%s",
		os.Getenv("ACI_CLIENT_ID"),
		os.Getenv("ACI_SECRET"),
	)
	return base64.StdEncoding.EncodeToString([]byte(authStr))
}

func buildBearer(tokenData tokenResponse) string {
	return fmt.Sprintf("Bearer %s", tokenData.AccessToken)
}

type tokenResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresOn   int64  `json:"expires_on"`
	IdToken     string `json:"id_token"`
	TokenType   string `json:"token_type"`
}

type clientResponse struct {
	TenantID string `json:"tenant_id"`
}

type usageItem struct {
	Tenant string                   `json:"tenant"`
	Usages []map[string]interface{} `json:"usages"`
}

type usageResponse struct {
	Items []usageItem `json:"items"`
}

type applicationResponse struct {
	Id   string `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"`
}

func main() {
	client := resty.New()
	client.SetBaseURL(fmt.Sprintf("%s/api/2", os.Getenv("ACI_DC_URL")))
	client.SetHeader("Accept", "application/json")

	// fetch token
	var tokenData tokenResponse
	resp, err := client.R().
		SetHeader("Authorization", fmt.Sprintf("Basic %s", encodeClientCredentials())).
		SetHeader("Content-Type", "application/x-www-form-urlencoded").
		SetFormData(map[string]string{"grant_type": "client_credentials"}).
		SetResult(&tokenData).
		Post("/idp/token")
	if err != nil {
		panic(err)
	}

	if !resp.IsSuccess() {
		fmt.Println("Unable to fetch token.")
		fmt.Printf("%v", string(resp.Body()))
		os.Exit(-1)
	}

	fmt.Printf("Got a token: %s***\n", tokenData.AccessToken[0:6])
	client.SetHeader("Authorization", buildBearer(tokenData))

	// fetch tenant id
	var clientData clientResponse
	resp, err = client.R().
		SetResult(&clientData).
		Get(fmt.Sprintf("/clients/%s", os.Getenv("ACI_CLIENT_ID")))
	if err != nil {
		panic(err)
	}
	if !resp.IsSuccess() {
		fmt.Println("Unable to fetch tenant id")
		os.Exit(-1)
	}

	fmt.Printf("Got tenant id: %s\n", clientData.TenantID)

	// fetch usage data
	var usageData usageResponse
	resp, err = client.R().
		SetQueryParams(map[string]string{"tenants": clientData.TenantID}).
		SetResult(&usageData).
		Get("/tenants/usages")
	if err != nil {
		panic(err)
	}
	if !resp.IsSuccess() {
		fmt.Println("Unable to fetch usage data.")
		os.Exit(-1)
	}

	for _, items := range usageData.Items {
		//fmt.Printf("Got tenant ID: %s\n", items.Tenant)
		for _, usages := range items.Usages {
			if usages["name"] != "hci_s3_storage" {
				continue
			}

			var applicationData applicationResponse
			resp, err = client.R().
				SetResult(&applicationData).
				Get(fmt.Sprintf("/applications/%s", usages["application_id"]))
			if !resp.IsSuccess() {
				panic(err)
			}

			fmt.Printf("%s (Type: %s)\n\n%s -- %.2f GB\n",
				applicationData.Name,
				applicationData.Type,
				usages["name"],
				// bitshift -> byte to gb
				(usages["absolute_value"].(float64) / (1 << 30)))
		}
	}
}

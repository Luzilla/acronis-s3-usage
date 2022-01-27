package acronis

import (
	"encoding/base64"
	"fmt"
	"os"

	"github.com/go-resty/resty/v2"
)

type AcronisClient struct {
	ClientID string
	Secret   string
	DCurl    string
	token    tokenResponse
	http     *resty.Client
}

func NewClient(clientID string, secret string, dcURL string) *AcronisClient {
	http := resty.New()
	http.SetBaseURL(fmt.Sprintf("%s/api/2", dcURL))
	http.SetHeader("Accept", "application/json")

	return &AcronisClient{
		ClientID: clientID,
		Secret:   secret,
		DCurl:    dcURL,
		http:     http,
	}
}

func (c *AcronisClient) GetApplication(appId string) (ApplicationResponse, error) {
	var data ApplicationResponse
	resp, err := c.http.R().
		SetResult(&data).
		Get(fmt.Sprintf("/applications/%s", appId))
	if !resp.IsSuccess() {
		return ApplicationResponse{}, err
	}

	return data, nil
}

// fetch tenant id
func (c *AcronisClient) GetTenantID() (string, error) {
	c.fetchToken()

	var data clientResponse
	resp, err := c.http.R().
		SetResult(&data).
		Get(fmt.Sprintf("/clients/%s", os.Getenv("ACI_CLIENT_ID")))
	if err != nil {
		return "", err
	}
	if !resp.IsSuccess() {
		return "", fmt.Errorf("unable to fetch tenant id: %s", resp.Body())
	}

	return data.TenantID, nil
}

// fetch usage data
func (c *AcronisClient) GetUsage(tenantId string) (UsageResponse, error) {
	c.fetchToken()

	var data UsageResponse
	resp, err := c.http.R().
		SetQueryParams(map[string]string{"tenants": tenantId}).
		SetResult(&data).
		Get("/tenants/usages")
	if err != nil {
		return UsageResponse{}, err
	}
	if !resp.IsSuccess() {
		return UsageResponse{}, fmt.Errorf("unable to fetch usage data: %s", resp.Body())
	}

	return data, nil
}

func (c *AcronisClient) encodeClientCredentials() string {
	authStr := fmt.Sprintf("%s:%s",
		c.ClientID,
		c.Secret,
	)
	return base64.StdEncoding.EncodeToString([]byte(authStr))
}

func (c *AcronisClient) buildBearer() string {
	return fmt.Sprintf("Bearer %s", c.token.AccessToken)
}

func (c *AcronisClient) fetchToken() {
	if c.token.AccessToken != "" {
		return
	}

	// fetch token
	var tokenData tokenResponse
	resp, err := c.http.R().
		SetHeader("Authorization", fmt.Sprintf("Basic %s", c.encodeClientCredentials())).
		SetHeader("Content-Type", "application/x-www-form-urlencoded").
		SetFormData(map[string]string{"grant_type": "client_credentials"}).
		SetResult(&tokenData).
		Post("/idp/token")
	if err != nil {
		panic(err)
	}

	if !resp.IsSuccess() { // FIXME
		fmt.Println("Unable to fetch token.")
		fmt.Printf("%v", string(resp.Body()))
		os.Exit(-1)
	}

	c.token = tokenData

	fmt.Printf("Got a token: %s***\n", c.token.AccessToken[0:8])
	c.http.SetHeader("Authorization", c.buildBearer())
}

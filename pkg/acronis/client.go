package acronis

import (
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/Luzilla/acronis-s3-usage/internal/client"
)

type AcronisClient struct {
	clientID string
	secret   string

	// data center URL
	dcURL string

	// token
	token tokenResponse

	// http client
	httpClient *client.Client
}

func NewClient(clientID string, secret string, dcURL string) *AcronisClient {
	httpClient := client.New(fmt.Sprintf("%s/api/2", dcURL), &http.Client{})

	return &AcronisClient{
		clientID:   clientID,
		secret:     secret,
		dcURL:      dcURL,
		httpClient: httpClient,
	}
}

func (c *AcronisClient) GetApplication(ctx context.Context, appId string) (ApplicationResponse, error) {
	var data ApplicationResponse

	req, err := c.authedRequest(ctx, http.MethodGet, fmt.Sprintf("/applications/%s", appId), nil)
	if err != nil {
		return data, err
	}

	if err := c.httpClient.DoInto(req, &data); err != nil {
		return data, err
	}

	return data, nil
}

// fetch tenant id
func (c *AcronisClient) GetTenantID(ctx context.Context) (string, error) {
	var data clientResponse

	req, err := c.authedRequest(ctx, http.MethodGet, fmt.Sprintf("/clients/%s", os.Getenv("ACI_CLIENT_ID")), nil)
	if err != nil {
		return "", err
	}

	if err := c.httpClient.DoInto(req, &data); err != nil {
		return "", err
	}

	return data.TenantID, nil
}

// fetch usage data
func (c *AcronisClient) GetUsage(ctx context.Context, tenantId string) (UsageResponse, error) {
	var data UsageResponse

	req, err := c.authedRequest(ctx, http.MethodGet, "/tenants/usages?tenants="+tenantId, nil)
	if err != nil {
		return data, err
	}

	if err := c.httpClient.DoInto(req, &data); err != nil {
		return data, err
	}

	return data, nil
}

func (c *AcronisClient) encodeClientCredentials() string {
	authStr := fmt.Sprintf("%s:%s",
		c.clientID,
		c.secret,
	)
	return base64.StdEncoding.EncodeToString([]byte(authStr))
}

func (c *AcronisClient) buildBearer() string {
	return fmt.Sprintf("Bearer %s", c.token.AccessToken)
}

// authedRequest creates a request with a valid bearer token.
func (c *AcronisClient) authedRequest(ctx context.Context, method, path string, body io.Reader) (*http.Request, error) {
	if err := c.fetchToken(ctx); err != nil {
		return nil, err
	}

	req, err := c.request(ctx, method, path, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", c.buildBearer())

	return req, nil
}

func (c *AcronisClient) fetchToken(ctx context.Context) error {
	if c.token.AccessToken != "" {
		return nil
	}

	form := url.Values{}
	form.Set("grant_type", "client_credentials")

	var tokenData tokenResponse

	req, err := c.request(ctx, http.MethodPost, "/idp/token", strings.NewReader(form.Encode()))
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Basic %s", c.encodeClientCredentials()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	if err := c.httpClient.DoInto(req, &tokenData); err != nil {
		return err
	}

	c.token = tokenData

	fmt.Printf("Got a token: %s***\n", c.token.AccessToken[0:5])
	return nil
}

func (c *AcronisClient) request(ctx context.Context, method, path string, body io.Reader) (*http.Request, error) {
	req, err := c.httpClient.NewRequest(ctx, method, path, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")
	return req, nil
}

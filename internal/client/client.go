package client

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Client struct {
	baseURL string
	http    *http.Client
}

func New(baseURL string, httpClient *http.Client) *Client {
	return &Client{
		baseURL: baseURL,
		http:    httpClient,
	}
}

// NewRequest creates an *http.Request with path appended to the client's baseURL.
// The path must start with "/" or "?" (e.g. "/applications/foo" or "/?ostor-users&key=val").
func (c *Client) NewRequest(ctx context.Context, method, path string, body io.Reader) (*http.Request, error) {
	return http.NewRequestWithContext(ctx, method, c.baseURL+path, body)
}

func (c *Client) Do(req *http.Request) (*http.Response, error) {
	return c.http.Do(req)
}

func (c *Client) DoInto(req *http.Request, resp any) error {
	res, err := c.http.Do(req)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	if res.StatusCode >= 400 {
		body, err := io.ReadAll(res.Body)
		if err != nil {
			return fmt.Errorf("unable to read the body while handling an error (status: %d): %w", res.StatusCode, err)
		}
		return fmt.Errorf("unable to complete request (status: %d): %s", res.StatusCode, body)
	}

	if err := json.NewDecoder(res.Body).Decode(resp); err != nil {
		return err
	}

	return nil
}

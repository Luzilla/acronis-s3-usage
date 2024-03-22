package ostor

import (
	"errors"
	"fmt"

	"github.com/go-resty/resty/v2"
)

type Ostor struct {
	client *resty.Client

	endpoint string

	// credentials (system user account)
	keyID, secretKeyID string
}

func New(endpoint, accessKeyID, secretKeyID string) (*Ostor, error) {
	http := resty.New()
	http.SetBaseURL(endpoint)
	http.SetHeader("Accept", "application/json")

	if len(endpoint) == 0 {
		return nil, errors.New("missing endpoint")
	}

	if len(accessKeyID) == 0 {
		return nil, errors.New("missing accessKeyID")
	}

	if len(secretKeyID) == 0 {
		return nil, errors.New("missing secretKeyID")
	}

	return &Ostor{
		client:      http,
		endpoint:    endpoint,
		keyID:       accessKeyID,
		secretKeyID: secretKeyID,
	}, nil
}

func (o *Ostor) getRequest(cmd, query string, into any) (*resty.Response, error) {
	signature, date, err := createSignature("GET", o.secretKeyID, cmd)
	if err != nil {
		return nil, fmt.Errorf("unable to create signature: %s", err)
	}

	o.client.Header.Set("Accept", "*/*")
	o.client.Header.Set("Date", date)
	o.client.Header.Set("Authorization", authHeader(o.keyID, signature))

	resp, err := o.client.R().SetResult(&into).Get("/?" + query)

	if err != nil {
		return nil, err
	}

	if resp.IsError() {
		headers := resp.Header()
		if headers.Get("X-Amz-Err-Message") != "" {
			return nil, fmt.Errorf("request failed: %s", headers.Get("X-Amz-Err-Message"))
		}

		return nil, fmt.Errorf("unable to make request: %d", resp.StatusCode())
	}

	return resp, nil
}

func (o *Ostor) postRequest(cmd, query string) (*resty.Response, error) {
	signature, date, err := createSignature("POST", o.secretKeyID, cmd)
	if err != nil {
		return nil, fmt.Errorf("unable to create signature: %s", err)
	}

	o.client.Header.Set("Accept", "*/*")
	o.client.Header.Set("Date", date)
	o.client.Header.Set("Authorization", authHeader(o.keyID, signature))

	return o.client.R().Post("/?" + query)
}

func (o *Ostor) putRequest(cmd, query string) (*resty.Response, error) {
	signature, date, err := createSignature("PUT", o.secretKeyID, cmd)
	if err != nil {
		return nil, fmt.Errorf("unable to create signature: %s", err)
	}

	o.client.Header.Set("Accept", "*/*")
	o.client.Header.Set("Date", date)
	o.client.Header.Set("Authorization", authHeader(o.keyID, signature))

	return o.client.R().Put("/?" + query)
}

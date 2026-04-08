package ostor

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sort"
	"strings"

	"github.com/Luzilla/acronis-s3-usage/internal/client"
)

var (
	// init
	errMissingEndpoint  = &OstorConfigError{msg: "missing endpoint"}
	errMissingAccessKey = &OstorConfigError{msg: "missing access key id"}
	errMissingSecretKey = &OstorConfigError{msg: "missing secret key id"}
)

type Ostor struct {
	client *client.Client

	endpoint string

	// credentials (system user account)
	keyID, secretKeyID string
}

func New(endpoint, accessKeyID, secretKeyID string) (*Ostor, error) {
	endpoint = strings.TrimSpace(endpoint)
	accessKeyID = strings.TrimSpace(accessKeyID)
	secretKeyID = strings.TrimSpace(secretKeyID)

	if len(endpoint) == 0 {
		return nil, errMissingEndpoint
	}

	if len(accessKeyID) == 0 {
		return nil, errMissingAccessKey
	}

	if len(secretKeyID) == 0 {
		return nil, errMissingSecretKey
	}

	c := client.New(endpoint, &http.Client{})

	return &Ostor{
		client:      c,
		endpoint:    endpoint,
		keyID:       accessKeyID,
		secretKeyID: secretKeyID,
	}, nil
}

func (o *Ostor) delete(ctx context.Context, cmd string, query map[string]string) (*http.Response, error) {
	return o.request(ctx, http.MethodDelete, cmd, query, nil)
}

func (o *Ostor) get(ctx context.Context, cmd string, query map[string]string, into any) (*http.Response, error) {
	return o.request(ctx, http.MethodGet, cmd, query, into)
}

func (o *Ostor) post(ctx context.Context, cmd string, query map[string]string, into any) (*http.Response, error) {
	return o.request(ctx, http.MethodPost, cmd, query, into)
}

func (o *Ostor) put(ctx context.Context, cmd string, query map[string]string, into any) (*http.Response, error) {
	return o.request(ctx, http.MethodPut, cmd, query, into)
}

func (o *Ostor) request(ctx context.Context, method, cmd string, query map[string]string, into any) (*http.Response, error) {
	signature, date, err := createSignature(method, o.secretKeyID, cmd)
	if err != nil {
		return nil, fmt.Errorf("unable to create signature: %w", err)
	}

	reqPath := buildPath(cmd, query)

	req, err := o.client.NewRequest(ctx, method, reqPath, nil)
	if err != nil {
		return nil, fmt.Errorf("unable to create request: %w", err)
	}

	req.Header.Set("Accept", "*/*")
	req.Header.Set("Date", date)
	req.Header.Set("Authorization", authHeader(o.keyID, signature))
	req.Header.Set("User-Agent", "Mozilla/5.0 (compatible; ostor-client/x.y; +https://github.com/Luzilla/acronis-s3-usage)")

	res, err := o.client.Do(req)
	if err != nil {
		return res, &OstorTransportError{
			Res: res,
			Err: err,
		}
	}

	// Buffer the body so callers can still read it
	bodyBytes, err := io.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		return res, &OstorTransportError{
			Res: res,
			Err: fmt.Errorf("unable to read response body: %w", err),
		}
	}
	res.Body = io.NopCloser(bytes.NewReader(bodyBytes))

	if res.StatusCode >= 400 {
		if msg := res.Header.Get("X-Amz-Err-Message"); msg != "" {
			return res, &OstorAPIError{
				Res: res,
				Err: errors.New(msg),
			}
		}

		return res, &OstorTransportError{
			Res: res,
			Err: fmt.Errorf("unable to make request: %d", res.StatusCode),
		}
	}

	if into != nil && len(bodyBytes) > 0 {
		if err := json.Unmarshal(bodyBytes, into); err != nil {
			return res, &OstorTransportError{
				Res: res,
				Err: fmt.Errorf("unable to decode response: %w", err),
			}
		}
	}

	return res, nil
}

// buildPath constructs the request path from the command and query parameters.
// The command becomes the first (valueless) query parameter, e.g. /?ostor-users&emailAddress=foo.
func buildPath(cmd string, query map[string]string) string {
	u := "/?" + cmd

	keys := make([]string, 0, len(query))
	for k := range query {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		v := query[k]
		if v == "" {
			u += "&" + url.QueryEscape(k)
		} else {
			u += "&" + url.QueryEscape(k) + "=" + url.QueryEscape(v)
		}
	}
	return u
}

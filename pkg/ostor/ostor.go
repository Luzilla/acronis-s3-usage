package ostor

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/go-resty/resty/v2"
)

var (
	// init
	errMissingEndpoint  = &OstorConfigError{msg: "missing endpoint"}
	errMissingAccessKey = &OstorConfigError{msg: "missing access key id"}
	errMissingSecretKey = &OstorConfigError{msg: "missing secret key id"}

	// usage
	errMethodNotSupported = &OstorUsageError{msg: "unsupported method"}
)

type Ostor struct {
	client *resty.Client

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

	client := resty.New()
	client.SetBaseURL(endpoint)

	return &Ostor{
		client:      client,
		endpoint:    endpoint,
		keyID:       accessKeyID,
		secretKeyID: secretKeyID,
	}, nil
}

func (o *Ostor) delete(cmd string, query map[string]string) (*http.Response, error) {
	return o.request(o.client.R().
		SetQueryParams(query), cmd, resty.MethodDelete, "/?"+cmd)
}

func (o *Ostor) get(cmd string, query map[string]string, into any) (*http.Response, error) {
	return o.request(o.client.R().
		SetQueryParams(query).
		SetResult(&into), cmd, resty.MethodGet, "/?"+cmd)
}

func (o *Ostor) post(cmd, query string, into any) (*http.Response, error) {
	request := o.client.R()
	if into != nil {
		request = request.SetResult(into)
	}
	return o.request(request, cmd, resty.MethodPost, "/?"+query)
}

func (o *Ostor) put(cmd, query string, into any) (*http.Response, error) {
	request := o.client.R()
	if into != nil {
		request = request.SetResult(&into)
	}
	return o.request(request, cmd, resty.MethodPut, "/?"+query)
}

func (o *Ostor) request(req *resty.Request, cmd, method, url string) (*http.Response, error) {
	signature, date, err := createSignature(method, o.secretKeyID, cmd)
	if err != nil {
		return nil, fmt.Errorf("unable to create signature: %s", err)
	}

	req.SetHeaderMultiValues(map[string][]string{
		http.CanonicalHeaderKey("Accept"):        {"*/*"},
		http.CanonicalHeaderKey("Date"):          {date},
		http.CanonicalHeaderKey("Authorization"): {authHeader(o.keyID, signature)},
		http.CanonicalHeaderKey("User-Agent"):    {"Mozilla/5.0 (compatible; ostor-client/x.y; +https://github.com/Luzilla/acronis-s3-usage)"},
	})

	var res *resty.Response

	switch method {
	case resty.MethodDelete:
		res, err = req.Delete(url)
	case resty.MethodGet:
		res, err = req.Get(url)
	case resty.MethodPost:
		res, err = req.Post(url)
	case resty.MethodPut:
		res, err = req.Put(url)
	default:
		// return early: this is a library problem
		return nil, errMethodNotSupported
	}

	if err != nil {
		return toHTTPResponse(res), &OstorTransportError{
			Res: toHTTPResponse(res),
			Err: err,
		}
	}

	httpRes := toHTTPResponse(res)

	if res.StatusCode() < 400 {
		return httpRes, nil
	}

	// error based on status code
	if res.Header().Get("X-Amz-Err-Message") != "" {
		return httpRes, &OstorAPIError{
			Res: httpRes,
			Err: errors.New(res.Header().Get("X-Amz-Err-Message")),
		}
	}

	return httpRes, &OstorTransportError{
		Res: httpRes,
		Err: fmt.Errorf("unable to make request: %d", res.StatusCode()),
	}
}

// toHTTPResponse converts a resty response to a stdlib *http.Response with
// the body replaced by a re-readable buffer (resty already consumed it).
func toHTTPResponse(res *resty.Response) *http.Response {
	if res == nil {
		return nil
	}

	raw := res.RawResponse
	raw.Body = io.NopCloser(bytes.NewReader(res.Body()))
	return raw
}

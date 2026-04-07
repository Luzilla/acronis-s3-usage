package ostor

import (
	"errors"
	"fmt"
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

func (o *Ostor) delete(cmd string, query map[string]string) (*resty.Response, error) {
	return o.request(o.client.R().
		SetQueryParams(query), cmd, resty.MethodDelete, "/?"+cmd)
}

func (o *Ostor) get(cmd string, query map[string]string, into any) (*resty.Response, error) {
	return o.request(o.client.R().
		SetQueryParams(query).
		SetResult(&into), cmd, resty.MethodGet, "/?"+cmd)
}

func (o *Ostor) post(cmd, query string, into any) (*resty.Response, error) {
	request := o.client.R()
	if into != nil {
		request = request.SetResult(into)
	}
	return o.request(request, cmd, resty.MethodPost, "/?"+query)
}

func (o *Ostor) put(cmd, query string, into any) (*resty.Response, error) {
	request := o.client.R()
	if into != nil {
		request = request.SetResult(&into)
	}
	return o.request(request, cmd, resty.MethodPut, "/?"+query)
}

func (o *Ostor) request(req *resty.Request, cmd, method, url string) (*resty.Response, error) {
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
		return res, errMethodNotSupported
	}

	if err != nil {
		// fmt.Printf("%v", res.Request)
		// b, _ := io.ReadAll(res.RawBody())
		// fmt.Println(b)
		return res, &OstorTransportError{
			Res: res,
			Err: err,
		}
	}

	if !res.IsError() {
		return res, nil
	}

	// error based on status code
	if res.Header().Get("X-Amz-Err-Message") != "" {
		return res, &OstorAPIError{
			Res: res,
			Err: errors.New(res.Header().Get("X-Amz-Err-Message")),
		}
	}

	return res, &OstorTransportError{
		Res: res,
		Err: fmt.Errorf("unable to make request: %d", res.StatusCode()),
	}
}

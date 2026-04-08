package ostor

import (
	"bytes"
	"errors"
	"io"
	"log/slog"
	"net/http"
)

type OstorConfigError struct {
	msg string
}

func (e *OstorConfigError) Error() string {
	return e.msg
}

type OstorAPIError struct {
	Res *http.Response
	Err error
}

func (e *OstorAPIError) Error() string {
	return e.Err.Error()
}

func (e *OstorAPIError) Unwrap() error {
	return e.Err
}

type OstorTransportError struct {
	Res *http.Response
	Err error
}

func (e *OstorTransportError) Error() string {
	return e.Err.Error()
}

func (e *OstorTransportError) Unwrap() error {
	return e.Err
}

func (e *OstorTransportError) LogValue() slog.Value {
	if e.Res == nil || e.Res.Request == nil {
		return slog.GroupValue(
			slog.String("error", e.Err.Error()),
		)
	}

	var body string
	if e.Res.Body != nil {
		b, err := io.ReadAll(e.Res.Body)
		if err == nil {
			body = string(b)
			// Reset the body so it remains readable for subsequent callers
			e.Res.Body = io.NopCloser(bytes.NewReader(b))
		}
	}

	return slog.GroupValue(
		slog.Group("request",
			slog.String("url", e.Res.Request.URL.String()),
			slog.String("method", e.Res.Request.Method),
			slog.String("signature", e.Res.Request.Header.Get("authorization")),
			slog.String("date", e.Res.Request.Header.Get("date")),
		),
		slog.Group("response",
			slog.Any("headers", e.Res.Header),
			slog.String("body", body),
		),
	)
}

// IsConfigError returns true when the client is mis-configured
func IsConfigError(err error) bool {
	var configErr *OstorConfigError
	return errors.As(err, &configErr)
}

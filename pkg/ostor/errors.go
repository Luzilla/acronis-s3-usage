package ostor

import (
	"errors"
	"log/slog"

	"github.com/go-resty/resty/v2"
)

type OstorConfigError struct {
	msg string
}

func (e *OstorConfigError) Error() string {
	return e.msg
}

type OstorUsageError struct {
	msg string
}

func (e *OstorUsageError) Error() string {
	return e.msg
}

type OstorAPIError struct {
	Res *resty.Response
	Err error
}

func (e *OstorAPIError) Error() string {
	return e.Err.Error()
}

type OstorTransportError struct {
	Res *resty.Response
	Err error
}

func (e *OstorTransportError) Error() string {
	return e.Err.Error()
}

func (e *OstorTransportError) LogValue() slog.Value {
	if e.Res == nil || e.Res.Request == nil {
		return slog.GroupValue(
			slog.String("error", e.Err.Error()),
		)
	}

	return slog.GroupValue(
		slog.Group("request",
			slog.String("url", e.Res.Request.URL),
			slog.String("method", e.Res.Request.Method),
			slog.String("signature", e.Res.Request.Header.Get("authorization")),
			slog.String("date", e.Res.Request.Header.Get("date")),
		),
		slog.Group("response",
			slog.Any("headers", e.Res.Header()),
			slog.String("body", string(e.Res.Body())),
		),
	)
}

// IsConfigError returns true when the client is mis-configured
func IsConfigError(err error) bool {
	var configErr *OstorConfigError
	return errors.As(err, &configErr)
}

// IsUsageError returns true when the library is doing something wrong
func IsUsageError(err error) bool {
	var usageErr *OstorUsageError
	return errors.As(err, &usageErr)
}

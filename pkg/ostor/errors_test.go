package ostor_test

import (
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/Luzilla/acronis-s3-usage/pkg/ostor"
	"github.com/go-resty/resty/v2"
)

func (s *OstorTestSuite) TestNew() {
	s.T().Run("no access key", func(t *testing.T) {
		_, err := ostor.New("http://blah", " ", "123")
		s.Require().Error(err)
		s.Assert().IsType(&ostor.OstorConfigError{}, err)
		s.Assert().True(ostor.IsConfigError(err))
	})
	s.T().Run("no secret key", func(t *testing.T) {
		_, err := ostor.New("http://blah", "123", "")
		s.Require().Error(err)
		s.Assert().IsType(&ostor.OstorConfigError{}, err)
		s.Assert().True(ostor.IsConfigError(err))
	})
	s.T().Run("no endpoint", func(t *testing.T) {
		_, err := ostor.New("", "123", "123")
		s.Require().Error(err)
		s.Assert().IsType(&ostor.OstorConfigError{}, err)
		s.Assert().True(ostor.IsConfigError(err))
	})
}

func (s *OstorTestSuite) TestTransportError() {
	header := http.Header{}
	header.Set("date", time.Now().Format(time.RFC1123Z))
	header.Set("authorization", "something")

	req := &resty.Request{
		URL:    "http://localhost",
		Method: "GET",
		Header: header,
	}

	body := strings.NewReader("fixture")

	res := &resty.Response{
		Request: req,
		RawResponse: &http.Response{
			StatusCode: 401,
			Body:       io.NopCloser(body),
		},
	}
	res.SetBody([]byte("fixture"))

	err := &ostor.OstorTransportError{
		Res: res,
		Err: fmt.Errorf("something"),
	}

	slog.Error("error message", "err", err)
}

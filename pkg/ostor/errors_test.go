package ostor_test

import (
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"testing"

	"github.com/Luzilla/acronis-s3-usage/pkg/ostor"
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
	reqURL, _ := url.Parse("http://localhost")

	req := &http.Request{
		URL:    reqURL,
		Method: "GET",
		Header: http.Header{
			"Date":          {s.T().Name()},
			"Authorization": {"something"},
		},
	}

	res := &http.Response{
		StatusCode: 401,
		Header:     http.Header{},
		Body:       http.NoBody,
		Request:    req,
	}

	err := &ostor.OstorTransportError{
		Res: res,
		Err: fmt.Errorf("something"),
	}

	slog.Error("error message", "err", err)
}

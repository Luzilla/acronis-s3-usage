package ostor_test

import (
	"log/slog"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/Luzilla/acronis-s3-usage/pkg/ostor"
	"github.com/Luzilla/acronis-s3-usage/pkg/ostormock"
	"github.com/stretchr/testify/suite"
)

func init() {
	// enable default logger
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})))
}

type OstorTestSuite struct {
	suite.Suite
	client     *ostor.Ostor
	mockServer *httptest.Server
}

func (s *OstorTestSuite) SetupSuite() {
	server, url := ostormock.StartMockServer(s.T())
	s.mockServer = server

	client, _ := ostor.New(url, "system", "system-pass")
	s.client = client
}

func (s *OstorTestSuite) TeardownSuite() {
	s.mockServer.CloseClientConnections()
	s.mockServer.Close()
}

func TestOstor(t *testing.T) {
	suite.Run(t, new(OstorTestSuite))
}

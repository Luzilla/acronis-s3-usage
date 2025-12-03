package ostor_test

import (
	"log/slog"
	"os"
	"testing"

	"github.com/Luzilla/acronis-s3-usage/pkg/ostor"
	"github.com/Luzilla/acronis-s3-usage/pkg/ostormock"
	"github.com/stretchr/testify/suite"

	"go.uber.org/goleak"
)

func init() {
	// enable default logger
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})))
}

type OstorTestSuite struct {
	suite.Suite
	client *ostor.Ostor
}

func (s *OstorTestSuite) SetupTest() {
	_, url := ostormock.StartMockServer(s.T())

	client, _ := ostor.New(url, "system", "system-pass")
	s.client = client
}

func TestMain(m *testing.M) {
	goleak.VerifyTestMain(m)
}

func TestOstor(t *testing.T) {
	suite.Run(t, new(OstorTestSuite))
}

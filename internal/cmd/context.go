package cmd

import (
	"context"

	"github.com/Luzilla/acronis-s3-usage/pkg/ostor"
)

func getOstorFromContext(ctx context.Context) *ostor.Ostor {
	return ctx.Value(ostorClient).(*ostor.Ostor)
}

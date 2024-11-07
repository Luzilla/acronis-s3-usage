package cmd

import (
	"context"

	"github.com/Luzilla/acronis-s3-usage/pkg/ostor"
	"github.com/urfave/cli/v2"
)

func Before(cCtx *cli.Context) error {
	client, err := ostor.New(
		cCtx.String("s3-endpoint"),
		cCtx.String("s3-system-key-id"),
		cCtx.String("s3-system-secret"))
	if err != nil {
		return err
	}

	cCtx.Context = context.WithValue(cCtx.Context, ostorClient, client)
	return nil
}

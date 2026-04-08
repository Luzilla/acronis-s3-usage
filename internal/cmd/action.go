package cmd

import (
	"context"
	"fmt"

	"github.com/Luzilla/acronis-s3-usage/pkg/ostor"
	"github.com/urfave/cli/v3"
)

func Before(ctx context.Context, c *cli.Command) (context.Context, error) {
	client, err := ostor.New(
		c.String("s3-endpoint"),
		c.String("s3-system-key-id"),
		c.String("s3-system-secret"))
	if err != nil {
		if ostor.IsConfigError(err) {
			return ctx, fmt.Errorf("please check your configuration (--s3-endpoint, --s3-system-key-id, --s3-system-secret): %w", err)
		}
		return ctx, err
	}

	ctx = context.WithValue(ctx, ostorClient, client)
	return ctx, nil
}

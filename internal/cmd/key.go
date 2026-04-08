package cmd

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/urfave/cli/v3"
)

func revokeKey(ctx context.Context, c *cli.Command) error {
	client := getOstorFromContext(ctx)

	email := c.String("email")
	keyID := c.String("key-id")

	_, err := client.RevokeKey(ctx, email, keyID)
	if err != nil {
		return err
	}
	slog.Info("success")

	return nil
}

func createKey(ctx context.Context, c *cli.Command) error {
	client := getOstorFromContext(ctx)

	email := c.String("email")

	_, _, err := client.GenerateCredentials(ctx, email)
	if err != nil {
		return err
	}
	slog.Info("success")

	return nil
}

func rotateKey(ctx context.Context, c *cli.Command) error {
	client := getOstorFromContext(ctx)

	email := c.String("email")
	keyID := c.String("key-id")

	keyPair, _, err := client.RotateKey(ctx, email, keyID)
	if err != nil {
		return err
	}

	fmt.Println("New key generated:")
	fmt.Printf("Access Key ID:     %s\n", keyPair.AccessKeyID)
	fmt.Printf("Secret Access Key: %s\n", keyPair.SecretAccessKey)
	fmt.Println("")

	return nil
}

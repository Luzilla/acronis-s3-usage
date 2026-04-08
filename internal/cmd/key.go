package cmd

import (
	"fmt"
	"log/slog"

	"github.com/urfave/cli/v2"
)

func revokeKey(cCtx *cli.Context) error {
	client := getOstorFromContext(cCtx.Context)

	email := cCtx.String("email")
	keyID := cCtx.String("key-id")

	_, err := client.RevokeKey(cCtx.Context, email, keyID)
	if err != nil {
		return err
	}
	slog.Info("success")

	return nil
}

func createKey(cCtx *cli.Context) error {
	client := getOstorFromContext(cCtx.Context)

	email := cCtx.String("email")

	_, _, err := client.GenerateCredentials(cCtx.Context, email)
	if err != nil {
		return err
	}
	slog.Info("success")

	return nil
}

func rotateKey(cCtx *cli.Context) error {
	client := getOstorFromContext(cCtx.Context)

	email := cCtx.String("email")
	keyID := cCtx.String("key-id")

	keyPair, _, err := client.RotateKey(cCtx.Context, email, keyID)
	if err != nil {
		return err
	}

	fmt.Println("New key generated:")
	fmt.Printf("Access Key ID:     %s\n", keyPair.AccessKeyID)
	fmt.Printf("Secret Access Key: %s\n", keyPair.SecretAccessKey)
	fmt.Println("")

	return nil
}

package cmd

import (
	"fmt"
	"log/slog"

	"github.com/Luzilla/acronis-s3-usage/internal/utils"
	"github.com/Luzilla/acronis-s3-usage/pkg/ostor"
	"github.com/rodaine/table"
	"github.com/urfave/cli/v2"
)

func Users(cCtx *cli.Context) error {
	client := cCtx.Context.Value(OstorClient).(*ostor.Ostor)

	users, _, err := client.ListUsers(cCtx.Bool("usage"))
	if err != nil {
		return err
	}

	var tbl table.Table
	if cCtx.Bool("usage") {
		tbl = table.New("Email", "ID", "State", "Space")
	} else {
		tbl = table.New("Email", "ID", "State")
	}

	tbl.WithHeaderFormatter(headerFmt()).WithFirstColumnFormatter(columnFmt())

	for _, u := range users.Users {
		if cCtx.Bool("usage") {
			tbl.AddRow(u.Email, u.ID, u.State, formatBytes(u.Space.Current))
		} else {
			tbl.AddRow(u.Email, u.ID, u.State)
		}
	}
	tbl.Print()

	return nil
}

func CreateUser(cCtx *cli.Context) error {
	client := cCtx.Context.Value(OstorClient).(*ostor.Ostor)

	email := cCtx.String("email")

	user, _, err := client.CreateUser(email)
	if err != nil {
		return err
	}

	fmt.Println("success")

	fmt.Printf("ID: %s\n", user.ID)
	fmt.Printf("Email: %s\n", user.Email)

	fmt.Println("Generated the following key-pair:")
	fmt.Printf("Key ID: %s\n", user.AccessKeys[0].AccessKeyID)
	fmt.Printf("Secret Access Key: %s\n", user.AccessKeys[0].SecretAccessKey)
	return nil
}

func ShowUser(cCtx *cli.Context) error {
	client := cCtx.Context.Value(OstorClient).(*ostor.Ostor)

	email := cCtx.String("email")

	user, resp, err := client.GetUser(email)
	if err != nil {
		if resp.StatusCode() == 404 {
			return fmt.Errorf("no user with email %q found", email)
		}
		return err
	}

	fmt.Printf("ID:    %s\n", user.ID)
	fmt.Printf("Email: %s\n", user.Email)
	fmt.Printf("State: %s\n", user.State)
	fmt.Println("")

	tblAK := table.New("Key ID", "Secret Key ID")
	tblAK.WithHeaderFormatter(headerFmt()).WithFirstColumnFormatter(columnFmt())

	for _, ak := range user.AccessKeys {
		tblAK.AddRow(ak.AccessKeyID, ak.SecretAccessKey)
	}

	tblAK.Print()

	fmt.Println("")

	buckets, _, err := client.GetBuckets(email)
	if err != nil {
		return err
	}

	tbl := table.New("Bucket", "Size (current)", "Created At")
	tbl.WithHeaderFormatter(headerFmt()).WithFirstColumnFormatter(columnFmt())

	for _, b := range buckets.Buckets {
		tbl.AddRow(b.Name, utils.PrettyByteSize(b.Size.Current), b.CreatedAt)
	}
	tbl.Print()

	return nil
}

func RevokeKey(cCtx *cli.Context) error {
	client := cCtx.Context.Value(OstorClient).(*ostor.Ostor)

	email := cCtx.String("email")
	keyID := cCtx.String("key-id")

	_, err := client.RevokeKey(email, keyID)
	if err != nil {
		return err
	}
	slog.Info("success")

	return nil
}

func CreateKey(cCtx *cli.Context) error {
	client := cCtx.Context.Value(OstorClient).(*ostor.Ostor)

	email := cCtx.String("email")

	_, _, err := client.GenerateCredentials(email)
	if err != nil {
		return err
	}
	slog.Info("success")

	return nil
}

func formatBytes(bytes int64) string {
	const (
		KB = 1024
		MB = KB * 1024
		GB = MB * 1024
		TB = GB * 1024
	)

	switch {
	case bytes >= TB:
		return fmt.Sprintf("%.2f TB", float64(bytes)/float64(TB))
	case bytes >= GB:
		return fmt.Sprintf("%.2f GB", float64(bytes)/float64(GB))
	case bytes >= MB:
		return fmt.Sprintf("%.2f MB", float64(bytes)/float64(MB))
	default:
		return fmt.Sprintf("%d Bytes", bytes)
	}
}

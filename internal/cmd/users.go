package cmd

import (
	"fmt"

	"github.com/Luzilla/acronis-s3-usage/internal/utils"
	"github.com/Luzilla/acronis-s3-usage/pkg/ostor"
	"github.com/rodaine/table"
	"github.com/urfave/cli/v2"
)

func Users(cCtx *cli.Context) error {
	client := cCtx.Context.Value(OstorClient).(*ostor.Ostor)

	resp, err := client.ListUsers()
	if err != nil {
		return err
	}

	tbl := table.New("Email", "ID", "State")
	tbl.WithHeaderFormatter(headerFmt()).WithFirstColumnFormatter(columnFmt())

	for _, u := range resp.Users {
		tbl.AddRow(u.Email, u.ID, u.State)
	}
	tbl.Print()

	return nil
}

func CreateUser(cCtx *cli.Context) error {
	client := cCtx.Context.Value(OstorClient).(*ostor.Ostor)

	email := cCtx.String("email")

	err := client.CreateUser(email)
	if err != nil {
		return err
	}

	fmt.Println("success")
	return nil
}

func ShowUser(cCtx *cli.Context) error {
	client := cCtx.Context.Value(OstorClient).(*ostor.Ostor)

	email := cCtx.String("email")

	user, err := client.GetUser(email)
	if err != nil {
		return err
	}

	fmt.Printf("Email: %s (State: %s)\n", user.Email, user.State)
	fmt.Println("")

	tblAK := table.New("Key ID", "Secret Key ID")
	tblAK.WithHeaderFormatter(headerFmt()).WithFirstColumnFormatter(columnFmt())

	for _, ak := range user.AccessKeys {
		tblAK.AddRow(ak.KeyID, ak.SecretKeyID)
	}

	tblAK.Print()

	fmt.Println("")

	buckets, err := client.GetBuckets(email)
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

	resp, err := client.RevokeKey(email, keyID)
	if err != nil {
		return err
	}
	fmt.Println("success")

	fmt.Printf("%s", resp)
	return nil
}

func CreateKey(cCtx *cli.Context) error {
	client := cCtx.Context.Value(OstorClient).(*ostor.Ostor)

	email := cCtx.String("email")

	_, err := client.GenerateCredentials(email)
	if err != nil {
		return err
	}
	fmt.Println("success")
	return nil
}

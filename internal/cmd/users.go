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
			tbl.AddRow(u.Email, u.ID, u.State, utils.PrettyByteSize(u.Space.Current))
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

func DeleteUser(cCtx *cli.Context) error {
	client := cCtx.Context.Value(OstorClient).(*ostor.Ostor)

	email := cCtx.String("email")

	resp, err := client.DeleteUser(email)
	if err != nil {
		fmt.Println(resp.Request.URL)

		return err
	}

	fmt.Println("Account deleted")
	return nil
}

func LockUser(cCtx *cli.Context) error {
	client := cCtx.Context.Value(OstorClient).(*ostor.Ostor)

	email := cCtx.String("email")
	err := lockUnLockUser(client, email, true)
	if err != nil {
		return err
	}

	fmt.Println("Locked the account")
	return nil
}

func UnlockUser(cCtx *cli.Context) error {
	client := cCtx.Context.Value(OstorClient).(*ostor.Ostor)

	email := cCtx.String("email")
	err := lockUnLockUser(client, email, false)
	if err != nil {
		return err
	}

	fmt.Println("Unlocked the account")
	return nil
}

func lockUnLockUser(client *ostor.Ostor, email string, lock bool) error {
	resp, err := client.LockUnlockUser(email, lock)
	fmt.Println(string(resp.Body()))
	return err
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

	if len(user.AccessKeys) > 0 {
		tblAK := table.New("Key ID", "Secret Key ID")
		tblAK.WithHeaderFormatter(headerFmt()).WithFirstColumnFormatter(columnFmt())

		for _, ak := range user.AccessKeys {
			tblAK.AddRow(ak.AccessKeyID, ak.SecretAccessKey)
		}

		tblAK.Print()

		fmt.Println("")
	} else {
		errorNoticeFmt("User does not have any keys.")
	}

	buckets, _, err := client.GetBuckets(email)
	if err != nil {
		return err
	}

	if len(buckets.Buckets) > 0 {
		tbl := table.New("Bucket", "Size (current)", "Created At")
		tbl.WithHeaderFormatter(headerFmt()).WithFirstColumnFormatter(columnFmt())

		for _, b := range buckets.Buckets {
			tbl.AddRow(b.Name, utils.PrettyByteSize(b.Size.Current), b.CreatedAt)
		}
		tbl.Print()
	} else {
		errorNoticeFmt("User does not have any buckets.")
	}

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

func RotateKey(cCtx *cli.Context) error {
	client := cCtx.Context.Value(OstorClient).(*ostor.Ostor)

	email := cCtx.String("email")
	keyID := cCtx.String("key-id")

	keyPair, _, err := client.RotateKey(email, keyID)
	if err != nil {
		return err
	}

	fmt.Println("New key generated:")
	fmt.Printf("Access Key ID:     %s\n", keyPair.AccessKeyID)
	fmt.Printf("Secret Access Key: %s\n", keyPair.SecretAccessKey)
	fmt.Println("")

	return nil
}

func UserLimits(cCtx *cli.Context) error {
	client := cCtx.Context.Value(OstorClient).(*ostor.Ostor)

	email := cCtx.String("email")

	limits, _, err := client.GetUserLimits(email)
	if err != nil {
		return nil
	}

	tbl := table.New("Limit", "Value")
	tbl.WithHeaderFormatter(headerFmt()).WithFirstColumnFormatter(columnFmt())

	tbl.AddRow("Ops Default (ops/sec)", limits.OpsDefault)
	tbl.AddRow("Ops List (ops/sec)", limits.OpsList)
	tbl.AddRow("Ops Delete (ops/sec)", limits.OpsDelete)
	tbl.AddRow("Ops Get (ops/sec)", limits.OpsGet)
	tbl.AddRow("Ops Put (ops/sec)", limits.OpsPut)
	tbl.AddRow("Bandwidth Out (kb/sec)", limits.BandwidthOut)

	tbl.Print()

	return nil
}

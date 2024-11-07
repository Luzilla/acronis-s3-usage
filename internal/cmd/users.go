package cmd

import (
	"fmt"

	"github.com/Luzilla/acronis-s3-usage/internal/utils"
	"github.com/Luzilla/acronis-s3-usage/pkg/ostor"
	"github.com/rodaine/table"
	"github.com/urfave/cli/v2"
)

func users(cCtx *cli.Context) error {
	client := cCtx.Context.Value(ostorClient).(*ostor.Ostor)

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

func createUser(cCtx *cli.Context) error {
	client := cCtx.Context.Value(ostorClient).(*ostor.Ostor)

	user, _, err := client.CreateUser(cCtx.String("email"))
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

func deleteUser(cCtx *cli.Context) error {
	client := cCtx.Context.Value(ostorClient).(*ostor.Ostor)

	resp, err := client.DeleteUser(cCtx.String("email"))
	if err != nil {
		fmt.Println(resp.Request.URL)

		return err
	}

	fmt.Println("Account deleted")
	return nil
}

func lockUser(cCtx *cli.Context) error {
	client := cCtx.Context.Value(ostorClient).(*ostor.Ostor)

	err := lockUnLockUser(client, cCtx.String("email"), true)
	if err != nil {
		return err
	}

	fmt.Println("Locked the account")
	return nil
}

func unlockUser(cCtx *cli.Context) error {
	client := cCtx.Context.Value(ostorClient).(*ostor.Ostor)

	err := lockUnLockUser(client, cCtx.String("email"), false)
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

func showUser(cCtx *cli.Context) error {
	client := cCtx.Context.Value(ostorClient).(*ostor.Ostor)

	user, resp, err := client.GetUser(cCtx.String("email"))
	if err != nil {
		if resp.StatusCode() == 404 {
			return fmt.Errorf("no user with email %q found", cCtx.String("email"))
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

	buckets, _, err := client.GetBuckets(cCtx.String("email"))
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

func userLimits(cCtx *cli.Context) error {
	client := cCtx.Context.Value(ostorClient).(*ostor.Ostor)

	limits, _, err := client.GetUserLimits(cCtx.String("email"))
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

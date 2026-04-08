package cmd

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"github.com/Luzilla/acronis-s3-usage/internal/utils"
	"github.com/Luzilla/acronis-s3-usage/pkg/ostor"
	"github.com/rodaine/table"
	"github.com/urfave/cli/v3"
)

func users(ctx context.Context, c *cli.Command) error {
	client := getOstorFromContext(ctx)

	users, _, err := client.ListUsers(ctx, c.Bool("usage"))
	if err != nil {
		return err
	}

	var tbl table.Table
	if c.Bool("usage") {
		tbl = table.New("Email", "ID", "State", "Space")
	} else {
		tbl = table.New("Email", "ID", "State")
	}

	tbl.WithHeaderFormatter(headerFmt()).WithFirstColumnFormatter(columnFmt())

	for _, u := range users.Users {
		if c.Bool("usage") {
			tbl.AddRow(u.Email, u.ID, u.State, utils.PrettyByteSize(u.Space.Current))
		} else {
			tbl.AddRow(u.Email, u.ID, u.State)
		}
	}
	tbl.Print()

	return nil
}

func createUser(ctx context.Context, c *cli.Command) error {
	client := getOstorFromContext(ctx)

	user, _, err := client.CreateUser(ctx, c.String("email"))
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

func deleteUser(ctx context.Context, c *cli.Command) error {
	client := getOstorFromContext(ctx)

	_, err := client.DeleteUser(ctx, c.String("email"))
	if err != nil {
		var transportErr *ostor.OstorTransportError
		if errors.As(err, &transportErr) {
			slog.Error("request failed", "err", transportErr)
		}
		return err
	}

	fmt.Println("Account deleted")
	return nil
}

func lockUser(ctx context.Context, c *cli.Command) error {
	client := getOstorFromContext(ctx)

	err := lockUnLockUser(ctx, client, c.String("email"), true)
	if err != nil {
		return err
	}

	fmt.Println("Locked the account")
	return nil
}

func unlockUser(ctx context.Context, c *cli.Command) error {
	client := getOstorFromContext(ctx)

	err := lockUnLockUser(ctx, client, c.String("email"), false)
	if err != nil {
		return err
	}

	fmt.Println("Unlocked the account")
	return nil
}

func lockUnLockUser(ctx context.Context, client *ostor.Ostor, email string, lock bool) error {
	resp, err := client.LockUnlockUser(ctx, email, lock)
	if err != nil {
		return err
	}
	body, _ := io.ReadAll(resp.Body)
	fmt.Println(string(body))
	return nil
}

func showUser(ctx context.Context, c *cli.Command) error {
	client := getOstorFromContext(ctx)

	user, _, err := client.GetUser(ctx, c.String("email"))
	if err != nil {
		var apiErr *ostor.OstorAPIError
		if errors.As(err, &apiErr) && apiErr.Res.StatusCode == http.StatusNotFound {
			return fmt.Errorf("no user with email %q found", c.String("email"))
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

	buckets, _, err := client.GetBuckets(ctx, c.String("email"))
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

func userLimits(ctx context.Context, c *cli.Command) error {
	client := getOstorFromContext(ctx)

	limits, _, err := client.GetUserLimits(ctx, c.String("email"))
	if err != nil {
		return err
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

package cmd

import (
	"context"
	"fmt"

	"github.com/Luzilla/acronis-s3-usage/internal/utils"
	"github.com/Luzilla/acronis-s3-usage/pkg/s3"
	"github.com/rodaine/table"
	"github.com/urfave/cli/v3"
)

// this executes the action on 'behalf' of the user by returning the account
// and using the first credential pair to run the delete operations
func deleteBucket(ctx context.Context, c *cli.Command) error {
	client := getOstorFromContext(ctx)

	s3, err := s3.NewS3(ctx, c.String("s3-endpoint"), c.String("email"), client)
	if err != nil {
		return err
	}

	bucketName := c.String("bucket")

	if _, err := s3.IsDeletable(ctx, bucketName); err != nil {
		return err
	}

	fmt.Println("Bucket " + bucketName + " can be deleted")

	if err := s3.DeleteBucket(ctx, bucketName); err != nil {
		return err
	}

	fmt.Println("Bucket deleted")
	return nil
}

func listBuckets(ctx context.Context, c *cli.Command) error {
	client := getOstorFromContext(ctx)

	buckets, _, err := client.GetBuckets(ctx, c.String("email"))
	if err != nil {
		return err
	}

	tbl := table.New("Bucket", "Size (current)", "Owner", "Created At")
	tbl.WithHeaderFormatter(headerFmt()).WithFirstColumnFormatter(columnFmt())

	if len(buckets.Buckets) > 0 {
		for _, b := range buckets.Buckets {
			tbl.AddRow(b.Name, utils.PrettyByteSize(b.Size.Current), b.OwnerID, b.CreatedAt)
		}
	} else {
		tbl.AddRow("no buckets")
	}

	tbl.Print()

	return nil
}

func showBucket(ctx context.Context, c *cli.Command) error {
	listBuckets(ctx, c) // display the filter view first

	client := getOstorFromContext(ctx)

	s3, err := s3.NewS3(ctx, c.String("s3-endpoint"), c.String("email"), client)
	if err != nil {
		return err
	}

	fmt.Println("")

	tbl := table.New("File", "Size")
	tbl.WithHeaderFormatter(headerFmt()).WithFirstColumnFormatter(columnFmt())

	for o := range s3.ListContents(ctx, c.String("bucket")) {
		tbl.AddRow(o.Key, utils.PrettyByteSize(o.Size), o.Owner.ID)
	}

	tbl.Print()

	return nil
}

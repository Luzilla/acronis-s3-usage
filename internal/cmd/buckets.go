package cmd

import (
	"context"
	"fmt"

	"github.com/Luzilla/acronis-s3-usage/internal/utils"
	"github.com/Luzilla/acronis-s3-usage/pkg/ostor"
	"github.com/Luzilla/acronis-s3-usage/pkg/s3"
	"github.com/rodaine/table"
	"github.com/urfave/cli/v2"
)

// this executes the action on 'behalf' of the user by returning the account
// and using the first credential pair to run the delete operations
func DeleteBucket(cCtx *cli.Context) error {
	client := cCtx.Context.Value(OstorClient).(*ostor.Ostor)

	s3, err := s3.NewS3(cCtx.String("s3-endpoint"), cCtx.String("email"), client)
	if err != nil {
		return err
	}

	ctx := context.Background()
	bucketName := cCtx.String("bucket")

	if _, err := s3.IsDeletable(ctx, bucketName); err != nil {
		return err
	}

	if err := s3.DeleteBucket(ctx, bucketName); err != nil {
		return err
	}

	fmt.Println("Bucket deleted")
	return nil
}

func ListBuckets(cCtx *cli.Context) error {
	client := cCtx.Context.Value(OstorClient).(*ostor.Ostor)

	buckets, _, err := client.GetBuckets(cCtx.String("email"))
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

func ShowBucket(cCtx *cli.Context) error {
	ListBuckets(cCtx) // display the filter view first

	client := cCtx.Context.Value(OstorClient).(*ostor.Ostor)

	s3, err := s3.NewS3(cCtx.String("s3-endpoint"), cCtx.String("email"), client)
	if err != nil {
		return err
	}

	fmt.Println("")

	tbl := table.New("File", "Size")
	tbl.WithHeaderFormatter(headerFmt()).WithFirstColumnFormatter(columnFmt())

	for o := range s3.ListContents(context.Background(), cCtx.String("bucket")) {
		tbl.AddRow(o.Key, utils.PrettyByteSize(o.Size), o.Owner.ID)
	}

	tbl.Print()

	return nil
}

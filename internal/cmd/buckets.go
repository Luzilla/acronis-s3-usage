package cmd

import (
	"fmt"

	"github.com/Luzilla/acronis-s3-usage/internal/utils"
	"github.com/Luzilla/acronis-s3-usage/pkg/s3"
	"github.com/rodaine/table"
	"github.com/urfave/cli/v2"
)

// this executes the action on 'behalf' of the user by returning the account
// and using the first credential pair to run the delete operations
func deleteBucket(cCtx *cli.Context) error {
	client := getOstorFromContext(cCtx.Context)

	s3, err := s3.NewS3(cCtx.Context, cCtx.String("s3-endpoint"), cCtx.String("email"), client)
	if err != nil {
		return err
	}

	bucketName := cCtx.String("bucket")

	if _, err := s3.IsDeletable(cCtx.Context, bucketName); err != nil {
		return err
	}

	fmt.Println("Bucket " + bucketName + " can be deleted")

	if err := s3.DeleteBucket(cCtx.Context, bucketName); err != nil {
		return err
	}

	fmt.Println("Bucket deleted")
	return nil
}

func listBuckets(cCtx *cli.Context) error {
	client := getOstorFromContext(cCtx.Context)

	buckets, _, err := client.GetBuckets(cCtx.Context, cCtx.String("email"))
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

func showBucket(cCtx *cli.Context) error {
	listBuckets(cCtx) // display the filter view first

	client := getOstorFromContext(cCtx.Context)

	s3, err := s3.NewS3(cCtx.Context, cCtx.String("s3-endpoint"), cCtx.String("email"), client)
	if err != nil {
		return err
	}

	fmt.Println("")

	tbl := table.New("File", "Size")
	tbl.WithHeaderFormatter(headerFmt()).WithFirstColumnFormatter(columnFmt())

	for o := range s3.ListContents(cCtx.Context, cCtx.String("bucket")) {
		tbl.AddRow(o.Key, utils.PrettyByteSize(o.Size), o.Owner.ID)
	}

	tbl.Print()

	return nil
}

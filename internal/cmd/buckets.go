package cmd

import (
	"github.com/Luzilla/acronis-s3-usage/internal/utils"
	"github.com/Luzilla/acronis-s3-usage/pkg/ostor"
	"github.com/rodaine/table"
	"github.com/urfave/cli/v2"
)

func ListBuckets(cCtx *cli.Context) error {
	client := cCtx.Context.Value(OstorClient).(*ostor.Ostor)

	buckets, err := client.GetBuckets("")
	if err != nil {
		return err
	}

	tbl := table.New("Bucket", "Size (current)", "Owner", "Created At")
	tbl.WithHeaderFormatter(headerFmt()).WithFirstColumnFormatter(columnFmt())

	for _, b := range buckets.Buckets {
		tbl.AddRow(b.Name, utils.PrettyByteSize(b.Size.Current), b.OwnerID, b.CreatedAt)
	}
	tbl.Print()

	return nil
}

package cmd

import (
	"fmt"
	"os"

	"github.com/Luzilla/acronis-s3-usage/internal/utils"
	"github.com/Luzilla/acronis-s3-usage/pkg/ostor"
	"github.com/rodaine/table"
	"github.com/urfave/cli/v2"
)

type stat struct {
	Put        int
	Get        int
	List       int
	Other      int
	Downloaded int
	Uploaded   int
}

func List(cCtx *cli.Context) error {
	client := cCtx.Context.Value(OstorClient).(*ostor.Ostor)

	items, err := client.List()
	if err != nil {
		return err
	}

	fmt.Printf("Found %d objects\n", items.Count)
	fmt.Printf("Truncated: %t\n", items.Truncated)

	if items.Count == 0 {
		fmt.Println("no items")
		return nil
	}

	keep := map[string]stat{}

	for _, obj := range items.Items {
		// fmt.Println(obj)
		usage, err := client.ObjectUsage(obj)
		if err != nil {
			fmt.Println("usage: " + err.Error())
			os.Exit(2)
		}

		for _, item := range usage.Items {
			var b = item.Key.Bucket
			if _, ok := keep[b]; !ok {
				keep[b] = stat{}
			}
			keep[b] = addToStruct(keep[b], item.Counters.Operations, item.Counters.Net)
		}
	}

	tbl := table.New("Bucket", "Put", "Get", "List", "Other", "Downloaded", "Uploaded")
	tbl.WithHeaderFormatter(headerFmt()).WithFirstColumnFormatter(columnFmt())

	for bucket, s := range keep {
		tbl.AddRow(bucket, s.Put, s.Get, s.List, s.Other, utils.PrettyByteSize(s.Downloaded), utils.PrettyByteSize(s.Uploaded))
	}

	tbl.Print()

	return nil
}

func addToStruct(k stat, o ostor.ItemCountersOps, n ostor.ItemCountersNet) stat {
	k.Put = k.Put + o.Put
	k.Get = k.Get + o.Get
	k.List = k.List + o.List
	k.Other = k.Other + o.Other
	k.Uploaded = k.Uploaded + n.Uploaded
	k.Downloaded = k.Downloaded + n.Downloaded

	return k
}

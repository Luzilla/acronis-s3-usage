package cmd

import (
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/Luzilla/acronis-s3-usage/internal/utils"
	"github.com/Luzilla/acronis-s3-usage/pkg/ostor"
	"github.com/rodaine/table"
	"github.com/urfave/cli/v2"
)

type stat struct {
	Put        int64
	Get        int64
	List       int64
	Other      int64
	Downloaded int64
	Uploaded   int64
}

// show stats is really expensive, it will (attempt to) crawl the entire `?ostor-usage` endpoint
// and look up individual entries when returned, so for n pages returned from `?ostor-usage`, it
// will make n * number of items returned requests to `?ostor-usage&obj=FOO`
func showStats(cCtx *cli.Context) error {
	client := cCtx.Context.Value(ostorClient).(*ostor.Ostor)

	// fmt.Printf("from: %s\n", cCtx.Timestamp("from").String())

	keep := map[string]stat{}
	var after *string

	// this is for pure display
	page := 0

	for {

		if after != nil {
			slog.Debug("next page", slog.String("after", *after))
			if strings.Contains(*after, "1970") {
				break
			}
		}

		page++
		items, _, err := client.List(after)
		if err != nil {
			return err
		}

		slog.Debug("found items", slog.Int("page", page), slog.Int("items", items.Count), slog.Bool("truncated", items.Truncated))

		if items.Count == 0 {
			break
		}

		for _, obj := range items.Items {
			// fmt.Println(obj)
			usage, _, err := client.ObjectUsage(obj)
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

			after = &obj
		}

		// end of the pagination - no more records
		if !items.Truncated {
			break
		}

	}

	if len(keep) == 0 {
		slog.Info("no stats")
		return nil
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

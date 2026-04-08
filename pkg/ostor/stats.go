package ostor

import (
	"context"
	"net/http"
)

const qStats string = "ostor-usage"

func (o *Ostor) List(ctx context.Context, after *string) (*OStorResponse, *http.Response, error) {
	var stats *OStorResponse

	queryString := map[string]string{}

	if after != nil {
		queryString["after"] = *after
	}

	// TODO(till):
	// - limit queries over 710 produce broken JSON response
	// - limit cannot be used with pagination (unless we implement other logic)
	// queryString["limit"] = strconv.Itoa(710)

	resp, err := o.get(ctx, qStats, queryString, &stats)
	return stats, resp, err
}

func (o *Ostor) ObjectUsage(ctx context.Context, object string) (*OStorObjectUsageResponse, *http.Response, error) {
	var usage *OStorObjectUsageResponse

	resp, err := o.get(ctx, qStats, map[string]string{"obj": object}, &usage)
	return usage, resp, err
}

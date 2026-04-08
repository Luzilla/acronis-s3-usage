package ostor

import (
	"context"
	"net/http"
)

const qLimits string = "ostor-limits"

func (o *Ostor) GetUserLimits(ctx context.Context, email string) (*OstorUserLimits, *http.Response, error) {
	var limits *OstorUserLimits
	resp, err := o.get(ctx, qLimits, emailMap(email), &limits)
	return limits, resp, err
}

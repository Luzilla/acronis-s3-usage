package ostor

import (
	"context"
	"net/http"
)

const qBuckets string = "ostor-buckets"

func (o *Ostor) GetBuckets(ctx context.Context, email string) (*OstorBucketListResponse, *http.Response, error) {
	var buckets *OstorBucketListResponse
	resp, err := o.get(ctx, qBuckets, emailMap(email), &buckets)
	return buckets, resp, err
}

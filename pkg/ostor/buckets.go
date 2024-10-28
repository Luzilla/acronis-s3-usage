package ostor

import "github.com/go-resty/resty/v2"

const qBuckets string = "ostor-buckets"

func (o *Ostor) GetBuckets(email string) (*OstorBucketListResponse, *resty.Response, error) {
	var buckets *OstorBucketListResponse
	resp, err := o.get(qBuckets, emailMap(email), &buckets)
	return buckets, resp, err
}

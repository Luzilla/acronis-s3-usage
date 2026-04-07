package ostor

import "net/http"

const qBuckets string = "ostor-buckets"

func (o *Ostor) GetBuckets(email string) (*OstorBucketListResponse, *http.Response, error) {
	var buckets *OstorBucketListResponse
	resp, err := o.get(qBuckets, emailMap(email), &buckets)
	return buckets, resp, err
}

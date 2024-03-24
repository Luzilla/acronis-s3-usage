package ostor

const qBuckets string = "ostor-buckets"

func (o *Ostor) GetBuckets(email string) (*OstorBucketListResponse, error) {
	var buckets *OstorBucketListResponse
	_, err := o.get(qBuckets, map[string]string{"emailAddress": email}, &buckets)
	return buckets, err
}

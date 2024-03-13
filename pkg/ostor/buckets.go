package ostor

const qBuckets string = "ostor-buckets"

func (o *Ostor) GetBuckets(email string) (*OstorBucketListResponse, error) {
	var buckets *OstorBucketListResponse
	_, err := o.getRequest(qBuckets, qBuckets+"&emailAddress="+email, &buckets)
	return buckets, err
}

package ostor

import "github.com/go-resty/resty/v2"

const qLimits string = "ostor-limits"

func (o *Ostor) GetUserLimits(email string) (*OstorUserLimits, *resty.Response, error) {
	var limits *OstorUserLimits
	resp, err := o.get(qLimits, emailMap(email), &limits)
	return limits, resp, err
}

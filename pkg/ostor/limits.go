package ostor

import "net/http"

const qLimits string = "ostor-limits"

func (o *Ostor) GetUserLimits(email string) (*OstorUserLimits, *http.Response, error) {
	var limits *OstorUserLimits
	resp, err := o.get(qLimits, emailMap(email), &limits)
	return limits, resp, err
}

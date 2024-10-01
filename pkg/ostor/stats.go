package ostor

import "github.com/go-resty/resty/v2"

const qStats string = "ostor-usage"

func (o *Ostor) List(after *string) (*OStorResponse, *resty.Response, error) {
	var stats *OStorResponse

	queryString := map[string]string{}

	if after != nil {
		queryString["after"] = *after
	}

	// TODO(till):
	// - limit queries over 710 produce broken JSON response
	// - limit cannot be used with pagination (unless we implement other logic)
	// queryString["limit"] = strconv.Itoa(710)

	resp, err := o.get(qStats, queryString, &stats)
	return stats, resp, err
}

func (o *Ostor) ObjectUsage(object string) (*OStorObjectUsageResponse, *resty.Response, error) {
	var usage *OStorObjectUsageResponse

	resp, err := o.get(qStats, map[string]string{"obj": object}, &usage)
	return usage, resp, err
}

package ostor

const qStats string = "ostor-usage"

func (o *Ostor) List(after *string) (*OStorResponse, error) {
	var stats *OStorResponse

	queryString := map[string]string{}

	if after != nil {
		queryString["after"] = *after
	}

	// TODO(till):
	// - limit queries over 710 produce broken JSON response
	// - limit cannot be used with pagination (unless we implement other logic)
	// queryString["limit"] = strconv.Itoa(710)

	_, err := o.get(qStats, queryString, &stats)
	if err != nil {
		return stats, err
	}

	return stats, err
}

func (o *Ostor) ObjectUsage(object string) (*OStorObjectUsageResponse, error) {
	var usage *OStorObjectUsageResponse

	_, err := o.get(qStats, map[string]string{"obj": object}, &usage)
	if err != nil {
		return usage, err
	}

	return usage, err
}

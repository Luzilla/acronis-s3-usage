package ostor

const qStats string = "ostor-usage"

func (o *Ostor) List() (*OStorResponse, error) {
	var stats *OStorResponse

	_, err := o.getRequest(qStats, qStats, &stats)
	if err != nil {
		return stats, err
	}

	return stats, err

}

func (o *Ostor) ObjectUsage(object string) (*OStorObjectUsageResponse, error) {
	var usage *OStorObjectUsageResponse

	_, err := o.getRequest(qStats, qStats+"&obj="+object, &usage)
	if err != nil {
		return usage, err
	}

	return usage, err
}

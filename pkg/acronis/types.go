package acronis

type ApplicationResponse struct {
	Id   string `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"`
}

type UsageItems struct {
	Tenant string      `json:"tenant"`
	Usages []UsageItem `json:"usages"`
}

type UsageItem struct {
	AbsoluteValue   float64     `json:"absolute_value"`
	ApplicationID   string      `json:"application_id"`
	Edition         interface{} `json:"edition"`
	InfraID         string      `json:"infra_id"`
	MeasurementUnit string      `json:"measurement_unit"`
	Name            string      `json:"name"`
	RangeStart      string      `json:"range_start"`
	TenantID        float64     `json:"tenant_id"`
	ItemType        string      `json:"type"`
	UsageName       string      `json:"usage_name"`
	Value           float64     `json:"value"`

	// not relevant but there
	OfferingItem struct {
		Status int `json:"status"`
		Quota  struct {
			Value   interface{} `json:"value"`
			Overage interface{} `json:"overage"`
			Version int         `json:"version"`
		} `json:"quota"`
	} `json:"offering_item,omitempty"`
}

type UsageResponse struct {
	Items []UsageItems `json:"items"`
}

type clientResponse struct {
	TenantID string `json:"tenant_id"`
}

type tokenResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresOn   int64  `json:"expires_on"`
	IdToken     string `json:"id_token"`
	TokenType   string `json:"token_type"`
}

package acronis

type ApplicationResponse struct {
	Id   string `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"`
}

type UsageItem struct {
	Tenant string                   `json:"tenant"`
	Usages []map[string]interface{} `json:"usages"`
}

type UsageResponse struct {
	Items []UsageItem `json:"items"`
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

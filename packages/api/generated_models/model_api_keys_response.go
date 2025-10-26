package generated_models

type ApiKeysResponse struct {

	ApiKeys []VaultApiKey `json:"apiKeys"`

	// Total number of API keys
	TotalCount int32 `json:"totalCount"`

	// Number of API keys per page
	PageSize int32 `json:"pageSize"`

	// Current page index (starting from 1)
	PageIndex int32 `json:"pageIndex"`
}

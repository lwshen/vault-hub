package generated_models

type VaultsResponse struct {

	// Page of vault records for the authenticated user
	Vaults []VaultLite `json:"vaults"`

	// Total number of vaults available for pagination
	TotalCount int32 `json:"totalCount"`

	// Number of vaults returned per page
	PageSize int32 `json:"pageSize"`

	// Current page index (starting from 1)
	PageIndex int32 `json:"pageIndex"`
}

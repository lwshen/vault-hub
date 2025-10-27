package models

type ConfigResponse struct {

	// Whether OIDC authentication is enabled
	OidcEnabled bool `json:"oidcEnabled"`

	// Whether transactional email is enabled
	EmailEnabled bool `json:"emailEnabled"`
}

package models

type EmailTokenResponse struct {

	// Indicates whether the email token request was accepted
	Success bool `json:"success"`

	// Machine-readable status code describing the outcome
	Code string `json:"code"`
}

package models

type StatusResponse struct {

	// Application version
	Version string `json:"version"`

	// Git commit hash
	Commit string `json:"commit"`

	// System operational status
	SystemStatus string `json:"systemStatus"`

	// Database connection status
	DatabaseStatus string `json:"databaseStatus"`
}

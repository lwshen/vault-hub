package models

type AuditMetricsResponse struct {

	// Total number of audit events in the last 30 days
	TotalEventsLast30Days int32 `json:"totalEventsLast30Days"`

	// Number of audit events in the last 24 hours
	EventsCountLast24Hours int32 `json:"eventsCountLast24Hours"`

	// Number of vault-related events in the last 30 days
	VaultEventsLast30Days int32 `json:"vaultEventsLast30Days"`

	// Number of API key-related events in the last 30 days
	ApiKeyEventsLast30Days int32 `json:"apiKeyEventsLast30Days"`
}

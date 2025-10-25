package models

type AuditLogsResponse struct {

	AuditLogs []AuditLog `json:"auditLogs"`

	// Total number of logs matching the filter criteria
	TotalCount int32 `json:"totalCount"`

	// Number of logs per page
	PageSize int32 `json:"pageSize"`

	// Current page index (starting from 0)
	PageIndex int32 `json:"pageIndex"`
}

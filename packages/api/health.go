package api

import "time"

// HealthCheck builds the response payload used by health endpoints.
func HealthCheck() HealthCheckResponse {
	status := "ok"
	now := time.Now()
	return HealthCheckResponse{
		Status:    &status,
		Timestamp: &now,
	}
}

package generated_models

import (
	"time"
)

type HealthCheckResponse struct {

	Status string `json:"status,omitempty"`

	Timestamp time.Time `json:"timestamp,omitempty"`
}

package model

// HealthStatus represents the structure of the health check response.
type HealthStatus struct {
	Database bool `json:"database"`
}

func NewHealthStatus(database bool) HealthStatus {
	return HealthStatus{
		Database: database,
	}
}

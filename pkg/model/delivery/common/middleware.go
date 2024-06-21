package common

import (
	"time"
)

// HTTPIncomingLog represents a log entry for incoming HTTP requests.
type HTTPIncomingLog struct {
	Location      string    `json:"location"`
	CorrelationID string    `json:"correlation_id"`
	Time          time.Time `json:"time"`
	RequestMethod string    `json:"request_method"`
	RequestURL    string    `json:"request_url"`
	ClientIP      string    `json:"client_ip"`
	UserAgent     string    `json:"user_agent"`
}

// HTTPOutgoingLog represents a log entry for outgoing HTTP responses.
type HTTPOutgoingLog struct {
	Location       string        `json:"location"`
	CorrelationID  string        `json:"correlation_id"`
	Time           time.Time     `json:"time"`
	RequestMethod  string        `json:"request_method"`
	RequestURL     string        `json:"request_url"`
	ClientIP       string        `json:"client_ip"`
	UserAgent      string        `json:"user_agent"`
	ResponseStatus int           `json:"response_status"`
	Duration       time.Duration `json:"duration"`
}

// NewHTTPIncomingLog creates a new instance of HTTPIncomingLog with the current time.
func NewHTTPIncomingLog(location, corelationID, method, url, ip, userAgent string) *HTTPIncomingLog {
	return &HTTPIncomingLog{
		Location:      location,
		CorrelationID: corelationID,
		Time:          time.Now(),
		RequestMethod: method,
		RequestURL:    url,
		ClientIP:      ip,
		UserAgent:     userAgent,
	}
}

// NewHTTPOutgoingLog creates a new instance of HTTPOutgoingLog with the current time and response details.
func NewHTTPOutgoingLog(location, corelationID, method, url, ip, userAgent string, status int, duration time.Duration) *HTTPOutgoingLog {
	return &HTTPOutgoingLog{
		Location:       location,
		CorrelationID:  corelationID,
		Time:           time.Now(),
		RequestMethod:  method,
		RequestURL:     url,
		ClientIP:       ip,
		UserAgent:      userAgent,
		ResponseStatus: status,
		Duration:       duration,
	}
}

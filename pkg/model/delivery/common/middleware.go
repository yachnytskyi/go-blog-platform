package common

import (
	"time"
)

// HTTPIncomingLog represents a log entry for incoming HTTP requests.
type HTTPIncomingLog struct {
	Location      string    `json:"location"`       // The location where the log entry is created.
	CorrelationID string    `json:"correlation_id"` // Unique identifier to correlate logs.
	Time          time.Time `json:"time"`           // The time when the log entry is created.
	RequestMethod string    `json:"request_method"` // The HTTP method of the incoming request.
	RequestURL    string    `json:"request_url"`    // The URL of the incoming request.
	ClientIP      string    `json:"client_ip"`      // The IP address of the client making the request.
	UserAgent     string    `json:"user_agent"`     // The user agent string from the client making the request.
}

// HTTPOutgoingLog represents a log entry for outgoing HTTP responses.
type HTTPOutgoingLog struct {
	Location       string        `json:"location"`        // The location where the log entry is created.
	CorrelationID  string        `json:"correlation_id"`  // Unique identifier to correlate logs.
	Time           time.Time     `json:"time"`            // The time when the log entry is created.
	RequestMethod  string        `json:"request_method"`  // The HTTP method of the request corresponding to the response.
	RequestURL     string        `json:"request_url"`     // The URL of the request corresponding to the response.
	ClientIP       string        `json:"client_ip"`       // The IP address of the client that made the request.
	UserAgent      string        `json:"user_agent"`      // The user agent string from the client that made the request.
	ResponseStatus int           `json:"response_status"` // The HTTP status code of the response.
	Duration       time.Duration `json:"duration"`        // The duration taken to process the request and send the response.
}

// NewHTTPIncomingLog creates a new instance of HTTPIncomingLog with the current time.
//
// Parameters:
// - location: The location where the log entry is created.
// - correlationID: Unique identifier to correlate logs.
// - method: The HTTP method of the incoming request.
// - url: The URL of the incoming request.
// - ip: The IP address of the client making the request.
// - userAgent: The user agent string from the client making the request.
//
// Returns:
// - A new instance of HTTPIncomingLog instance populated with the provided values and the current time.
func NewHTTPIncomingLog(location, correlationID, method, url, ip, userAgent string) HTTPIncomingLog {
	return HTTPIncomingLog{
		Location:      location,
		CorrelationID: correlationID,
		Time:          time.Now(),
		RequestMethod: method,
		RequestURL:    url,
		ClientIP:      ip,
		UserAgent:     userAgent,
	}
}

// NewHTTPOutgoingLog creates a new instance of HTTPOutgoingLog with the current time and response details.
//
// Parameters:
// - location: The location where the log entry is created.
// - correlationID: Unique identifier to correlate logs.
// - method: The HTTP method of the request corresponding to the response.
// - url: The URL of the request corresponding to the response.
// - ip: The IP address of the client that made the request.
// - userAgent: The user agent string from the client that made the request.
// - status: The HTTP status code of the response.
// - duration: The duration taken to process the request and send the response.
//
// Returns:
// - A new instance of HTTPOutgoingLog instance populated with the provided values and the current time.
func NewHTTPOutgoingLog(location, correlationID, method, url, ip, userAgent string, status int, duration time.Duration) HTTPOutgoingLog {
	return HTTPOutgoingLog{
		Location:       location,
		CorrelationID:  correlationID,
		Time:           time.Now(),
		RequestMethod:  method,
		RequestURL:     url,
		ClientIP:       ip,
		UserAgent:      userAgent,
		ResponseStatus: status,
		Duration:       duration,
	}
}

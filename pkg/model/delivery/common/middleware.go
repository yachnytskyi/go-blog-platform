package common

import "time"

// HTTPLog represents a log entry for an HTTP request.
type HTTPLog struct {
	Location       string        `json:"location"`        // Location in the code where the log was generated.
	Time           time.Time     `json:"time"`            // The time when the log entry was created.
	RequestMethod  string        `json:"request_method"`  // The HTTP method used for the request (e.g., GET, POST).
	RequestURL     string        `json:"request_url"`     // The URL of the request.
	ClientIP       string        `json:"client_ip"`       // The IP address of the client making the request.
	UserAgent      string        `json:"user_agent"`      // The User-Agent header from the request.
	ResponseStatus int           `json:"response_status"` // The HTTP status code of the response.
	Duration       time.Duration `json:"duration"`        // The duration of the request handling.
}

// NewHTTPLog creates a new instance of HTTPLog with the current time.
func NewHTTPLog(location string, method, url, ip, userAgent string, status int, duration time.Duration) *HTTPLog {
	return &HTTPLog{
		Location:       location,
		Time:           time.Now(),
		RequestMethod:  method,
		RequestURL:     url,
		ClientIP:       ip,
		UserAgent:      userAgent,
		ResponseStatus: status,
		Duration:       duration,
	}
}

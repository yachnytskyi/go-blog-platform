package common

import "time"

// HTTPLog represents a log entry for an HTTP request.
type HTTPLog struct {
	Location       string        // Location in the code where the log was generated.
	Time           time.Time     // The time when the log entry was created.
	RequestMethod  string        // The HTTP method used for the request (e.g., GET, POST).
	RequestURL     string        // The URL of the request.
	ClientIP       string        // The IP address of the client making the request.
	UserAgent      string        // The User-Agent header from the request.
	ResponseStatus int           // The HTTP status code of the response.
	Duration       time.Duration // The duration of the request handling.
}

package common

import (
	"fmt"
	"time"

	constants "github.com/yachnytskyi/golang-mongo-grpc/config/constants"
)

type HTTPIncomingLog struct {
	Location      string    `json:"location"`       // The location where the log entry is created.
	CorrelationID string    `json:"correlation_id"` // Unique identifier to correlate logs.
	Time          time.Time `json:"time"`           // The time when the log entry is created.
	RequestMethod string    `json:"request_method"` // The HTTP method of the incoming request.
	RequestURL    string    `json:"request_url"`    // The URL of the incoming request.
	ClientIP      string    `json:"client_ip"`      // The IP address of the client making the request.
	UserAgent     string    `json:"user_agent"`     // The user agent string from the client making the request.
}

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

func (log HTTPIncomingLog) Error() string {
	return fmt.Sprintf(
		"location: %s, "+
			"correlation_id: %s, "+
			"time: %s, "+
			"request_method: %s, "+
			"request_url: %s, "+
			"client_ip: %s, "+
			"user_agent: %s",
		log.Location,
		log.CorrelationID,
		log.Time.Format(constants.DateTimeFormat),
		log.RequestMethod,
		log.RequestURL,
		log.ClientIP,
		log.UserAgent,
	)
}

func (log HTTPOutgoingLog) Error() string {
	return fmt.Sprintf(
		"location: %s, "+
			"correlation_id: %s, "+
			"time: %s, "+
			"request_method: %s, "+
			"request_url: %s, "+
			"client_ip: %s, "+
			"user_agent: %s, "+
			"response_status: %d, "+
			"duration: %s",
		log.Location,
		log.CorrelationID,
		log.Time.Format(constants.DateTimeFormat),
		log.RequestMethod,
		log.RequestURL,
		log.ClientIP,
		log.UserAgent,
		log.ResponseStatus,
		log.Duration,
	)
}

package common

import "time"

type HTTPLog struct {
	Location       string
	Time           time.Time
	RequestMethod  string
	RequestURL     string
	ClientIP       string
	UserAgent      string
	ResponseStatus int
	Duration       time.Duration
}

package common

import (
	"time"

	constants "github.com/yachnytskyi/golang-mongo-grpc/config/constants"
)

// FormatDate formats a time.Time instance to a custom string format.
func FormatDate(data time.Time) string {
	return data.Format(constants.DateTimeFormat)
}

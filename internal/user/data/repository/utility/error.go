package utility

import (
	"go.mongodb.org/mongo-driver/mongo"
)

func IsMongoDBError(err error) bool {
	if mongo.IsNetworkError(err) || mongo.IsTimeout(err) {
		return true
	}

	return false
}

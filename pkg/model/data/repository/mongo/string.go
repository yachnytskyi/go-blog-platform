package mongo

import (
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// BSONToString maps a bson.M query to a simple string representation.
func BSONToString(query bson.M) string {
	var result string
	for key, value := range query {
		switch valueType := value.(type) {
		case primitive.ObjectID:
			result += fmt.Sprintf("%s: %s, ", key, valueType.Hex())
		default:
			result += fmt.Sprintf("%s: %s, ", key, value)
		}
	}
	if len(result) > 0 {
		result = result[:len(result)-2] // Remove trailing comma and space
	}
	return result
}

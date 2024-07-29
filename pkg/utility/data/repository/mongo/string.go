package mongo

import (
	"fmt"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// BSONToString maps a bson.M query to a simple string representation.
func BSONToStringMapper(query bson.M) string {
	var builder strings.Builder

	for key, value := range query {
		switch valueType := value.(type) {
		case primitive.ObjectID:
			builder.WriteString(fmt.Sprintf("%s: %s, ", key, valueType.Hex()))
		case string:
			builder.WriteString(fmt.Sprintf("%s: %s, ", key, valueType))
		default:
			builder.WriteString(fmt.Sprintf("%s: %v, ", key, value))
		}
	}

	result := builder.String()
	if result != "" {
		result = result[:len(result)-2] // Remove trailing comma and space
	}

	return result
}

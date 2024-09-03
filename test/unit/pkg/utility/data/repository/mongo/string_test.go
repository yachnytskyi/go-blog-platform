package mongo

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	utility "github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/data/repository/mongo"
	test "github.com/yachnytskyi/golang-mongo-grpc/test"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Helper function to split result into key-value pairs
func parseResult(result string) map[string]string {
	pairs := strings.Split(result, ", ")
	parsed := make(map[string]string, len(pairs))
	for _, pair := range pairs {
		parts := strings.SplitN(pair, ": ", 2)
		if len(parts) == 2 {
			parsed[parts[0]] = parts[1]
		}
	}
	return parsed
}

func TestBSONToStringMapperNestedDocumentsAndArrays(t *testing.T) {
	t.Parallel()
	objectID := primitive.NewObjectID()
	embeddedDoc := bson.M{"subfield1": "subvalue1"}
	nestedArray := []interface{}{embeddedDoc, "stringValue", 42}

	query := bson.M{
		"field1": objectID,
		"field2": bson.M{
			"nestedField1": nestedArray,
			"nestedField2": 3.14,
		},
	}
	expected := map[string]string{
		"field1": objectID.Hex(),
		"field2": "map[nestedField1:[map[subfield1:subvalue1] stringValue 42] nestedField2:3.14]",
	}

	result := utility.BSONToStringMapper(query)
	parsedResult := parseResult(result)
	assert.Equal(t, expected, parsedResult, test.EqualMessage)
}

func TestBSONToStringMapperEmptyQuery(t *testing.T) {
	t.Parallel()
	query := bson.M{}
	expected := map[string]string{}

	result := utility.BSONToStringMapper(query)
	parsedResult := parseResult(result)
	assert.Equal(t, expected, parsedResult, test.EqualMessage)
}

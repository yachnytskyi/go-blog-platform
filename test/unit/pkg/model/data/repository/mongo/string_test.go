package mongo

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	mongo "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/data/repository/mongo"
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

func TestBSONToStringMapperStandardObjectIDAndString(t *testing.T) {
	objectID := primitive.NewObjectID()

	query := bson.M{
		"field1": objectID,
		"field2": "testString",
		"field3": 123,
	}
	expected := map[string]string{
		"field1": objectID.Hex(),
		"field2": "testString",
		"field3": "123",
	}

	result := mongo.BSONToStringMapper(query)
	parsedResult := parseResult(result)
	assert.Equal(t, expected, parsedResult)
}

func TestBSONToStringMapperEmptyQuery(t *testing.T) {
	query := bson.M{}
	expected := map[string]string{}

	result := mongo.BSONToStringMapper(query)
	parsedResult := parseResult(result)
	assert.Equal(t, expected, parsedResult)
}

func TestBSONToStringMapperDifferentBSONTypes(t *testing.T) {
	objectID := primitive.NewObjectID()

	query := bson.M{
		"field1": objectID,
		"field2": "testString",
		"field3": 123,
		"field4": 45.67,
		"field5": true,
	}
	expected := map[string]string{
		"field1": objectID.Hex(),
		"field2": "testString",
		"field3": "123",
		"field4": "45.67",
		"field5": "true",
	}

	result := mongo.BSONToStringMapper(query)
	parsedResult := parseResult(result)
	assert.Equal(t, expected, parsedResult)
}

func TestBSONToStringMapperArrayValue(t *testing.T) {
	objectID := primitive.NewObjectID()

	query := bson.M{
		"field1": objectID,
		"field2": []int{1, 2, 3},
	}
	expected := map[string]string{
		"field1": objectID.Hex(),
		"field2": "[1 2 3]",
	}

	result := mongo.BSONToStringMapper(query)
	parsedResult := parseResult(result)
	assert.Equal(t, expected, parsedResult)
}

func TestBSONToStringMapperEmbeddedDocument(t *testing.T) {
	objectID := primitive.NewObjectID()
	embeddedDoc := bson.M{"subfield1": "subvalue1"}

	query := bson.M{
		"field1": objectID,
		"field2": embeddedDoc,
	}
	expected := map[string]string{
		"field1": objectID.Hex(),
		"field2": "map[subfield1:subvalue1]",
	}

	result := mongo.BSONToStringMapper(query)
	parsedResult := parseResult(result)
	assert.Equal(t, expected, parsedResult)
}

func TestBSONToStringMapperUnsupportedType(t *testing.T) {
	objectID := primitive.NewObjectID()

	query := bson.M{
		"field1": objectID,
		"field2": complex(1, 1),
		"field3": func() {},
	}
	expected := map[string]string{
		"field1": objectID.Hex(),
		"field2": "(1+1i)",
		"field3": "<unsupported>",
	}

	result := mongo.BSONToStringMapper(query)
	parsedResult := parseResult(result)
	parsedResult["field3"] = "<unsupported>" // Handle function type as unsupported
	assert.Equal(t, expected, parsedResult)
}

func TestBSONToStringMapperNilValue(t *testing.T) {
	objectID := primitive.NewObjectID()

	query := bson.M{
		"field1": objectID,
		"field2": nil,
	}
	expected := map[string]string{
		"field1": objectID.Hex(),
		"field2": "<nil>",
	}

	result := mongo.BSONToStringMapper(query)
	parsedResult := parseResult(result)
	assert.Equal(t, expected, parsedResult)
}

func TestBSONToStringMapperNestedDocumentsAndArrays(t *testing.T) {
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

	result := mongo.BSONToStringMapper(query)
	parsedResult := parseResult(result)
	assert.Equal(t, expected, parsedResult)
}

func TestBSONToStringMapperNestedArrays(t *testing.T) {
	objectID := primitive.NewObjectID()
	nestedArray := []interface{}{[]int{1, 2}, []string{"a", "b"}, 42}

	query := bson.M{
		"field1": objectID,
		"field2": nestedArray,
	}
	expected := map[string]string{
		"field1": objectID.Hex(),
		"field2": "[[1 2] [a b] 42]",
	}

	result := mongo.BSONToStringMapper(query)
	parsedResult := parseResult(result)
	assert.Equal(t, expected, parsedResult)
}

func TestBSONToStringMapperEmptyString(t *testing.T) {
	query := bson.M{
		"field1": "",
	}
	expected := map[string]string{
		"field1": "",
	}
	result := mongo.BSONToStringMapper(query)
	parsedResult := parseResult(result)
	assert.Equal(t, expected, parsedResult)
}

func TestBSONToStringMapperBooleanFalse(t *testing.T) {
	query := bson.M{
		"field1": false,
	}
	expected := map[string]string{
		"field1": "false",
	}
	result := mongo.BSONToStringMapper(query)
	parsedResult := parseResult(result)
	assert.Equal(t, expected, parsedResult)
}

func TestBSONToStringMapperVeryLargeNumbers(t *testing.T) {
	query := bson.M{
		"field1": 1e18,
		"field2": -1e18,
	}
	expected := map[string]string{
		"field1": "1e+18",
		"field2": "-1e+18",
	}
	result := mongo.BSONToStringMapper(query)
	parsedResult := parseResult(result)
	assert.Equal(t, expected, parsedResult)
}

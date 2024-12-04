package mongo

import (
	"testing"

	"github.com/stretchr/testify/assert"
	model "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/data/repository/mongo"
	domain "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/error/domain"
	validator "github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/validator"
	test "github.com/yachnytskyi/golang-mongo-grpc/test"
	mock "github.com/yachnytskyi/golang-mongo-grpc/test/unit/mock/common"
	"go.mongodb.org/mongo-driver/bson"
)

const (
	location                  = "test.unit.pkg.mode.data.repository.mongo."
	emptyObjectID             = "000000000000000000000000"
	emptyHexString            = "Id is empty"
	invalidHexString          = "the provided hex string is not a valid ObjectID"
	invalidBsonMarshalMessage = "no encoder found for chan int"
	invalidBsonUnmarshalMsg   = "bson.Unmarshal failed"
)

// Tests for HexToObjectIDMapper

func TestHexToObjectIDMapperValidHex(t *testing.T) {
	t.Parallel()
	mockLogger := mock.NewMockLogger()

	validHex := "507f191e810c19729de860ea"
	result := model.HexToObjectIDMapper(mockLogger, location+"TestHexToObjectIDMapperValidHex", validHex)

	assert.False(t, validator.IsError(result.Error), test.NotFailureMessage)
	assert.Equal(t, validHex, result.Data.Hex(), test.EqualMessage)
}

func TestDataToMongoDocumentMapperSuccess(t *testing.T) {
	t.Parallel()
	mockLogger := mock.NewMockLogger()

	key := "name"
	value := "test"
	incomingData := map[string]interface{}{key: value}
	result := model.DataToMongoDocumentMapper(mockLogger, location+"TestDataToMongoDocumentMapperSuccess", incomingData)
	expectedDocument := bson.D{{Key: key, Value: value}}

	assert.False(t, validator.IsError(result.Error), test.ErrorNilMessage)
	assert.Equal(t, expectedDocument, *result.Data, test.EqualMessage)
}

func TestHexToObjectIDMapperEmptyHex(t *testing.T) {
	t.Parallel()
	mockLogger := mock.NewMockLogger()

	emptyHex := ""
	result := model.HexToObjectIDMapper(mockLogger, location+"TestHexToObjectIDMapperEmptyHex", emptyHex)
	expectedLocation := location + "TestHexToObjectIDMapperEmptyHex.HexToObjectIDMapper"
	expectedError := domain.NewInternalError(expectedLocation, emptyHexString)

	assert.True(t, validator.IsError(result.Error), test.ErrorNotNilMessage)
	assert.Equal(t, emptyObjectID, result.Data.Hex(), test.EqualMessage)
	assert.IsType(t, domain.InternalError{}, result.Error, test.EqualMessage)
	assert.Equal(t, expectedError, result.Error, test.EqualMessage)
}

func TestHexToObjectIDMapperInvalidHex(t *testing.T) {
	t.Parallel()
	mockLogger := mock.NewMockLogger()

	invalidHex := "12345"
	result := model.HexToObjectIDMapper(mockLogger, location+"TestHexToObjectIDMapperInvalidHex", invalidHex)
	expectedLocation := location + "TestHexToObjectIDMapperInvalidHex.HexToObjectIDMapper"
	expectedError := domain.NewInternalError(expectedLocation+".primitive.ObjectIDFromHex", invalidHexString)

	assert.True(t, validator.IsError(result.Error), test.ErrorNotNilMessage)
	assert.Equal(t, emptyObjectID, result.Data.Hex(), test.EqualMessage)
	assert.IsType(t, domain.InternalError{}, result.Error, test.EqualMessage)
	assert.Equal(t, expectedError, result.Error, test.EqualMessage)
}

func TestDataToMongoDocumentMapperBsonMarshalError(t *testing.T) {
	t.Parallel()
	mockLogger := mock.NewMockLogger()

	incomingData := make(chan int)
	result := model.DataToMongoDocumentMapper(mockLogger, location+"TestDataToMongoDocumentMapperBsonMarshalError", incomingData)
	expectedLocation := location + "TestDataToMongoDocumentMapperBsonMarshalError.DataToMongoDocumentMapper.bson.Marshal"
	expectedError := domain.NewInternalError(expectedLocation, invalidBsonMarshalMessage)

	assert.Nil(t, result.Data, test.DataNilMessage)
	assert.True(t, validator.IsError(result.Error), test.ErrorNotNilMessage)
	assert.IsType(t, domain.InternalError{}, result.Error, test.EqualMessage)
	assert.Equal(t, expectedError, result.Error, test.EqualMessage)
}

package mongo

import (
	commonModel "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/common"
	domainError "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/error/domain"
	logging "github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/logging"
	validator "github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/validator"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	emptyID = "Id is empty"
)

// DataToMongoDocumentMapper maps the incoming data to a BSON document.
// It performs the following steps:
// 1. Marshals the incoming data to BSON format.
// 2. Checks for errors during marshaling and logs them if any.
// 3. Unmarshals the BSON data into a BSON document.
// 4. Checks for errors during unmarshaling and logs them if any.
// 5. Returns the BSON document wrapped in a commonModel.Result.
//
// Parameters:
// - location (string): A string representing the location or context for error logging.
// - incomingData (any): The data to be mapped to a BSON document.
//
// Returns:
// - commonModel.Result[*bson.D]: The result containing either the BSON document or an error.
func DataToMongoDocumentMapper(location string, incomingData any) commonModel.Result[*bson.D] {
	data, err := bson.Marshal(incomingData)
	if validator.IsError(err) {
		internalError := domainError.NewInternalError(location+".DataToMongoDocumentMapper.bson.Marshal", err.Error())
		logging.Logger(internalError)
		return commonModel.NewResultOnFailure[*bson.D](internalError)
	}

	var document bson.D
	err = bson.Unmarshal(data, &document)
	if err != nil {
		internalError := domainError.NewInternalError(location+".DataToMongoDocumentMapper.bson.Unmarshal", err.Error())
		logging.Logger(internalError)
		return commonModel.NewResultOnFailure[*bson.D](internalError)
	}

	return commonModel.NewResultOnSuccess[*bson.D](&document)
}

// HexToObjectIDMapper maps a hexadecimal string representation of MongoDB ObjectID to its corresponding primitive.ObjectID type.
// It performs the following steps:
// 1. Checks if the provided ID is empty and logs an error if it is.
// 2. Attempts to map the hexadecimal string to a primitive.ObjectID.
// 3. Checks for errors during mapping and logs them if any.
// 4. Returns the mapped ObjectID wrapped in a commonModel.Result.
//
// Parameters:
// - location (string): A string representing the location or context for error logging.
// - id (string): The hexadecimal string to be mapped.
//
// Returns:
// - commonModel.Result[primitive.ObjectID]: The result containing either the mapped ObjectID or an error.
func HexToObjectIDMapper(location, id string) commonModel.Result[primitive.ObjectID] {
	if len(id) == 0 {
		internalError := domainError.NewInternalError(location+".HexToObjectIDMapper", emptyID)
		logging.Logger(internalError)
		return commonModel.NewResultOnFailure[primitive.ObjectID](internalError)
	}

	objectID, objectIDFromHexError := primitive.ObjectIDFromHex(id)
	if validator.IsError(objectIDFromHexError) {
		internalError := domainError.NewInternalError(location+".HexToObjectIDMapper.primitive.ObjectIDFromHex", objectIDFromHexError.Error())
		logging.Logger(internalError)
		return commonModel.NewResultOnFailure[primitive.ObjectID](internalError)
	}

	return commonModel.NewResultOnSuccess[primitive.ObjectID](objectID)
}

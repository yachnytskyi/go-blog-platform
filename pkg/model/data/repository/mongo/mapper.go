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

// DataToMongoDocumentMapper converts the incoming data to a BSON document.
// It uses BSON marshaling and unmarshaling to perform the conversion.
// Parameters:
// - ctx: A context to support cancellation and timeout.
// - location: A string representing the location or context for error logging.
// - incomingData: The data to be converted to a BSON document.
// Returns:
// - A Result containing a pointer to a BSON document.
// - A Result containing an error if the conversion fails.
func DataToMongoDocumentMapper(location string, incomingData any) commonModel.Result[*bson.D] {
	// Marshal incoming data to BSON format.
	data, err := bson.Marshal(incomingData)
	if validator.IsError(err) {
		// Log and handle the marshaling error.
		internalError := domainError.NewInternalError(location+".DataToMongoDocumentMapper.bson.Marshal", err.Error())
		logging.Logger(internalError)
		return commonModel.NewResultOnFailure[*bson.D](internalError)
	}

	// Unmarshal the BSON data into a BSON document.
	var document bson.D
	err = bson.Unmarshal(data, &document)
	if err != nil {
		// Log and handle the unmarshaling error.
		internalError := domainError.NewInternalError(location+".DataToMongoDocumentMapper.bson.Unmarshal", err.Error())
		logging.Logger(internalError)
		return commonModel.NewResultOnFailure[*bson.D](internalError)
	}

	return commonModel.NewResultOnSuccess[*bson.D](&document)
}

// HexToObjectIDMapper converts a hexadecimal string representation of MongoDB ObjectID
// to its corresponding primitive.ObjectID type.
// Parameters:
// - location: A string representing the location or context for error logging.
// - id: The hexadecimal string to be converted.
// Returns:
// - A Result containing the converted ObjectID if successful.
// - A Result containing an error if the conversion fails.
func HexToObjectIDMapper(location, id string) commonModel.Result[primitive.ObjectID] {
	// Check if the provided id is empty.
	if id == "" {
		// Log and handle the empty id error.
		internalError := domainError.NewInternalError(location+".HexToObjectIDMapper", emptyID)
		logging.Logger(internalError)
		// Return a failure result with the internal error.
		return commonModel.NewResultOnFailure[primitive.ObjectID](internalError)
	}

	// Convert the hexadecimal string to primitive.ObjectID.
	objectID, objectIDFromHexError := primitive.ObjectIDFromHex(id)
	if validator.IsError(objectIDFromHexError) {
		// Log and handle the conversion error.
		internalError := domainError.NewInternalError(location+".HexToObjectIDMapper.primitive.ObjectIDFromHex", objectIDFromHexError.Error())
		logging.Logger(internalError)
		// Return a default ObjectID and the error.
		return commonModel.NewResultOnFailure[primitive.ObjectID](internalError)
	}

	// Return the successfully converted ObjectID.
	return commonModel.NewResultOnSuccess[primitive.ObjectID](objectID)
}

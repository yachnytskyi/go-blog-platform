package mongo

import (
	domainError "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/error/domain"
	logging "github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/logging"
	validator "github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/validator"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	location = "pkg.model.data.repository.mongo."
)

// DataToMongoDocument converts the incoming data to a BSON document.
// It uses BSON marshaling and unmarshaling to perform the conversion.
func DataToMongoDocumentMapper(incomingData any) (document *bson.D, err error) {
	// Marshal incoming data to BSON format.
	data, err := bson.Marshal(incomingData)
	if validator.IsErrorNotNil(err) {
		// Log and handle the marshaling error.
		internalError := domainError.NewInternalError(location+"MongoMapper.bson.Marshal", err.Error())
		logging.Logger(internalError)
		return document, err
	}

	// Unmarshal the BSON data into a BSON document.
	err = bson.Unmarshal(data, &document)
	if validator.IsErrorNotNil(err) {
		// Log and handle the unmarshaling error.
		internalError := domainError.NewInternalError(location+"MongoMapper.bson.UnMarshal", err.Error())
		logging.Logger(internalError)
		return document, err
	}
	return
}

// HexToObjectID converts a hexadecimal string representation of MongoDB ObjectID
// to its corresponding primitive.ObjectID type.
// It takes a location string for context in error messages and the id as a string.
// Returns the converted ObjectID or an error if the conversion fails.
func HexToObjectIDMapper(location, id string) (primitive.ObjectID, error) {
	// Convert the hexadecimal string to primitive.ObjectID.
	objectID, objectIDFromHexError := primitive.ObjectIDFromHex(id)
	if validator.IsErrorNotNil(objectIDFromHexError) {
		// If an error occurs, create an internal error with context and log it.
		internalError := domainError.NewInternalError(location+".HexToObjectID", objectIDFromHexError.Error())
		logging.Logger(internalError)
		// Return a default ObjectID and the error.
		return primitive.NilObjectID, internalError
	}
	// Return the successfully converted ObjectID and nil error.
	return objectID, nil
}

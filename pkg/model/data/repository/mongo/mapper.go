package mongo

import (
	interfaces "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/interfaces"
	common "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/common"
	domain "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/error/domain"
	validator "github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/validator"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	emptyID = "Id is empty"
)

func DataToMongoDocumentMapper(logger interfaces.Logger, location string, incomingData any) common.Result[*bson.D] {
	data, err := bson.Marshal(incomingData)
	if validator.IsError(err) {
		internalError := domain.NewInternalError(location+".DataToMongoDocumentMapper.bson.Marshal", err.Error())
		logger.Error(internalError)
		return common.NewResultOnFailure[*bson.D](internalError)
	}

	var document bson.D
	err = bson.Unmarshal(data, &document)
	if validator.IsError(err) {
		internalError := domain.NewInternalError(location+".DataToMongoDocumentMapper.bson.Unmarshal", err.Error())
		logger.Error(internalError)
		return common.NewResultOnFailure[*bson.D](internalError)
	}

	return common.NewResultOnSuccess[*bson.D](&document)
}

func HexToObjectIDMapper(logger interfaces.Logger, location, id string) common.Result[primitive.ObjectID] {
	if id == "" {
		internalError := domain.NewInternalError(location+".HexToObjectIDMapper", emptyID)
		logger.Error(internalError)
		return common.NewResultOnFailure[primitive.ObjectID](internalError)
	}

	objectID, objectIDFromHexError := primitive.ObjectIDFromHex(id)
	if validator.IsError(objectIDFromHexError) {
		internalError := domain.NewInternalError(location+".HexToObjectIDMapper.primitive.ObjectIDFromHex", objectIDFromHexError.Error())
		logger.Error(internalError)
		return common.NewResultOnFailure[primitive.ObjectID](internalError)
	}

	return common.NewResultOnSuccess[primitive.ObjectID](objectID)
}

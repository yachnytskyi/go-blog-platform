package mongo

import (
	model "github.com/yachnytskyi/golang-mongo-grpc/pkg/dependency/model"
	common "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/common"
	domainError "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/error/domain"
	validator "github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/validator"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	emptyID = "Id is empty"
)

func DataToMongoDocumentMapper(logger model.Logger, location string, incomingData any) common.Result[*bson.D] {
	data, err := bson.Marshal(incomingData)
	if validator.IsError(err) {
		internalError := domainError.NewInternalError(location+".DataToMongoDocumentMapper.bson.Marshal", err.Error())
		logger.Error(internalError)
		return common.NewResultOnFailure[*bson.D](internalError)
	}

	var document bson.D
	err = bson.Unmarshal(data, &document)
	if err != nil {
		internalError := domainError.NewInternalError(location+".DataToMongoDocumentMapper.bson.Unmarshal", err.Error())
		logger.Error(internalError)
		return common.NewResultOnFailure[*bson.D](internalError)
	}

	return common.NewResultOnSuccess[*bson.D](&document)
}

func HexToObjectIDMapper(logger model.Logger, location, id string) common.Result[primitive.ObjectID] {
	if len(id) == 0 {
		internalError := domainError.NewInternalError(location+".HexToObjectIDMapper", emptyID)
		logger.Error(internalError)
		return common.NewResultOnFailure[primitive.ObjectID](internalError)
	}

	objectID, objectIDFromHexError := primitive.ObjectIDFromHex(id)
	if validator.IsError(objectIDFromHexError) {
		internalError := domainError.NewInternalError(location+".HexToObjectIDMapper.primitive.ObjectIDFromHex", objectIDFromHexError.Error())
		logger.Error(internalError)
		return common.NewResultOnFailure[primitive.ObjectID](internalError)
	}

	return common.NewResultOnSuccess[primitive.ObjectID](objectID)
}

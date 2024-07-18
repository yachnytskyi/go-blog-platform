package mongo

import (
	commonModel "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/common"
	domainError "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/error/domain"
	logger "github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/logger"
	validator "github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/validator"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	emptyID = "Id is empty"
)

func DataToMongoDocumentMapper(location string, incomingData any) commonModel.Result[*bson.D] {
	data, err := bson.Marshal(incomingData)
	if validator.IsError(err) {
		internalError := domainError.NewInternalError(location+".DataToMongoDocumentMapper.bson.Marshal", err.Error())
		logger.Logger(internalError)
		return commonModel.NewResultOnFailure[*bson.D](internalError)
	}

	var document bson.D
	err = bson.Unmarshal(data, &document)
	if err != nil {
		internalError := domainError.NewInternalError(location+".DataToMongoDocumentMapper.bson.Unmarshal", err.Error())
		logger.Logger(internalError)
		return commonModel.NewResultOnFailure[*bson.D](internalError)
	}

	return commonModel.NewResultOnSuccess[*bson.D](&document)
}

func HexToObjectIDMapper(location, id string) commonModel.Result[primitive.ObjectID] {
	if len(id) == 0 {
		internalError := domainError.NewInternalError(location+".HexToObjectIDMapper", emptyID)
		logger.Logger(internalError)
		return commonModel.NewResultOnFailure[primitive.ObjectID](internalError)
	}

	objectID, objectIDFromHexError := primitive.ObjectIDFromHex(id)
	if validator.IsError(objectIDFromHexError) {
		internalError := domainError.NewInternalError(location+".HexToObjectIDMapper.primitive.ObjectIDFromHex", objectIDFromHexError.Error())
		logger.Logger(internalError)
		return commonModel.NewResultOnFailure[primitive.ObjectID](internalError)
	}

	return commonModel.NewResultOnSuccess[primitive.ObjectID](objectID)
}

package mongo

import (
	domainError "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/error/domain"
	"github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/logging"
	"go.mongodb.org/mongo-driver/bson"
)

const (
	location = "pkg.model.data.repository.mongo."
)

// MongoMapper converts the incoming data to a BSON document.
// It uses BSON marshaling and unmarshaling to perform the conversion.
func MongoMapper(incomingData any) (document *bson.D, err error) {
	// Marshal incoming data to BSON format.
	data, err := bson.Marshal(incomingData)
	if err != nil {
		// Log and handle the marshaling error.
		internalError := domainError.NewInternalError(location+"MongoMapper.bson.Marshal", err.Error())
		logging.Logger(internalError)
		return document, err
	}

	// Unmarshal the BSON data into a BSON document.
	err = bson.Unmarshal(data, &document)
	if err != nil {
		// Log and handle the unmarshaling error.
		internalError := domainError.NewInternalError(location+"MongoMapper.bson.UnMarshal", err.Error())
		logging.Logger(internalError)
		return document, err
	}
	return
}

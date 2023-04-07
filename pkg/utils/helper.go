package utils

import "go.mongodb.org/mongo-driver/bson"

func MongoMapping(incomingData interface{}) (document *bson.D, err error) {
	data, err := bson.Marshal(incomingData)

	if err != nil {
		return document, err
	}

	err = bson.Unmarshal(data, &document)

	if err != nil {
		return document, err
	}

	return
}

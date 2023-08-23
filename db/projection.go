package db

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var CommonProtoctedFields []string = []string{"password", "device", "metadata.token_version", "metadata.provider"}

func ProtectFields(fields []string) primitive.D {
	var ExcludeFields primitive.D
	for _, field := range fields {
		ExcludeFields = append(ExcludeFields, bson.E{Key: field, Value: 0})
	}
	return ExcludeFields
}

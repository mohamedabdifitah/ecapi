package utils

import (
	"errors"
	"fmt"
	"strings"

	"go.mongodb.org/mongo-driver/mongo"
)

type ErrorResponse struct {
	Status  int
	Message interface{}
}

func HandlerError(err error) ErrorResponse {
	if mongo.IsTimeout(err) {
		return ErrorResponse{
			Status:  503, // or even you can 408 which means request timeout
			Message: fmt.Errorf("Request took too long to respond"),
		}
	}
	if mongo.IsNetworkError(err) {
		return ErrorResponse{
			Status:  500,
			Message: err,
		}
	}
	var e mongo.WriteException
	if errors.As(err, &e) {
		for _, we := range e.WriteErrors {
			switch we.Code {
			// duplicate filed
			case 11000:
				var field string = strings.Split(strings.Split(strings.Split(we.Message, "dup key: {")[1], "}")[0], ":")[0]
				var value string = strings.Split(strings.Split(strings.Split(strings.Split(we.Message, "dup key: {")[1], "}")[0], "\"")[0], "\"")[0]
				var reason map[string]interface{} = make(map[string]interface{})
				var me map[string]interface{} = make(map[string]interface{})
				me["reason"] = "duplicate"
				me["field"] = field
				me["message"] = fmt.Sprintf("resource: %v with value: %v already exists", field, value)
				reason["error"] = me
				return ErrorResponse{
					Status:  409,
					Message: reason, // fmt.Sprintf("resource: %v with value: %v already exists", field, value)
				}
			// Can't extract geo keys
			case 16755:
				return ErrorResponse{
					Status:  400,
					Message: err,
				}

			}
		}
	}
	return ErrorResponse{
		Status:  500,
		Message: "internal server error , please try again",
	}
}
func HandleRequestErrors() {
}

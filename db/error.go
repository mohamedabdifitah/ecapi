package db

import (
	"errors"
	"fmt"
	"strings"

	"go.mongodb.org/mongo-driver/mongo"
)

func IsDup(err error) bool {
	var e mongo.WriteException
	if errors.As(err, &e) {
		for _, we := range e.WriteErrors {
			if we.Code == 11000 {
				strings.Split(we.Message, "dup key: {")
				fmt.Println(we.Message)
				return true
			}
		}
	}
	return false
}
func DBErrorHandler(err error) *ErrorResponse {

	if mongo.IsTimeout(err) {
		return &ErrorResponse{
			Status:  503, // or even you can 408 which means request timeout
			Message: fmt.Errorf("request took too long to respond"),
			Type:    "string",
		}
	}
	if mongo.IsNetworkError(err) {
		return &ErrorResponse{
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
				// var field string = strings.Split(strings.Split(we.Message, "dup key: {")[1], "}")[0]

				return &ErrorResponse{
					Status:  409,
					Message: err,
				}
			// Can't extract geo keys
			case 16755:
				return &ErrorResponse{
					Status:  400,
					Message: err,
				}

			}
		}
	}
	return &ErrorResponse{
		Status:  500,
		Message: err,
	}
}

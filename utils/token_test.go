package utils

import (
	"testing"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

var token string

func TestGenerateAccessToken(t *testing.T) {
	var err error
	token, err = GenerateAccessToken("tes@emial.com", primitive.NewObjectID(), "role")
	t.Error(err)
}
func TestVerifyAccessToken(t *testing.T) {

	VerifyAccessToken(token)
	// eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2OTA5NjU5OTAsImlkIjoiT2JqZWN0SUQoXCI2NGM4YTA5NzQxZGY5ZjZkMTYwY2M1NDZcIikiLCJlbWFpbCI6ImV4YW1wbGVAZG9tYWluLmNvbSIsInJvbGUiOiIifQ.PNKKrOAw8ftQ_aHf-ZOqzmPe0SJ52C0O9dCOmSmcqFQ")
}

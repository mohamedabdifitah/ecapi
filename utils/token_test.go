package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestGenerateAccessTokenVerify(t *testing.T) {
	var err error
	token, err := GenerateAccessToken("tes@emial.com", primitive.NewObjectID(), "role")
	if err != nil {
		t.Error(err)
	}
	claim, err := VerifyAccessToken(token)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, claim.Email, "tes@emial.com")
}

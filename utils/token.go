package utils

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AcessTokenClaim struct {
	jwt.StandardClaims
	Id    string `json:"id"`
	Email string `json:"email"`
	Role string `json:"role"`
}
type RefereshTokenClaim struct {
	jwt.StandardClaims
	Id           string `json:"id"`
	Email        string `json:"email"`
	TokenVersion int    `json:"token_v"`
}

func GenerateAccessToken(email string, id primitive.ObjectID,role string) (string, error) {
	secretKey := os.Getenv("TOKEN_ACCESS_SECRET")
	claims := AcessTokenClaim{
		StandardClaims: jwt.StandardClaims{
			// set token lifetime in timestamp
			ExpiresAt: time.Now().Add(time.Duration(1500)).Unix(),
		},
		// add custom claims like user_id or email,
		// it can vary according to requirements

		Id:    id.String(),
		Email: email,
	}
	// generate a string using claims and HS256 algorithm
	tokenString := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// sign the generated key using secretKey
	token, err := tokenString.SignedString([]byte(secretKey))

	return token, err
}
func GenerateRefreshToken(email string, id primitive.ObjectID, v int) (string, error) {
	secretKey := os.Getenv("TOKEN_REF_SECRET")
	claims := RefereshTokenClaim{
		StandardClaims: jwt.StandardClaims{
			// set token lifetime in timestamp
			ExpiresAt: time.Now().Add(time.Duration(1500)).Unix(),
		},
		// add custom claims like user_id or email,
		// it can vary according to requirements

		Id:           id.String(),
		Email:        email,
		TokenVersion: v,
	}
	// generate a string using claims and HS256 algorithm
	tokenString := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// sign the generated key using secretKey
	token, err := tokenString.SignedString([]byte(secretKey))

	return token, err
}

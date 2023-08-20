package utils

import (
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AcessTokenClaim struct {
	jwt.StandardClaims
	Id    string `json:"id"`
	Email string `json:"email"`
	Role  string `json:"role"`
}
type RefereshTokenClaim struct {
	jwt.StandardClaims
	Id           string `json:"id"`
	Email        string `json:"email"`
	TokenVersion int    `json:"token_v"`
}

var (
	accessSecretKey  []byte = []byte(os.Getenv("TOKEN_ACCESS_SECRET"))
	refreshSecretKey []byte = []byte(os.Getenv("TOKEN_REF_SECRET"))
)

func GenerateAccessToken(email string, id primitive.ObjectID, role string) (string, error) {
	// secretKey := os.Getenv("TOKEN_ACCESS_SECRET")
	claims := AcessTokenClaim{
		StandardClaims: jwt.StandardClaims{
			// set token lifetime in timestamp
			ExpiresAt: time.Now().Add(time.Minute * 15).Unix(),
		},
		// add custom claims like user_id or email,
		// it can vary according to requirements

		Id:    id.Hex(),
		Email: email,
		Role:  role,
	}
	// generate a string using claims and HS256 algorithm
	tokenString := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// sign the generated key using secretKey
	token, err := tokenString.SignedString(accessSecretKey)

	return token, err
}
func GenerateRefreshToken(email string, id primitive.ObjectID, v int) (string, error) {
	// secretKey := os.Getenv("TOKEN_REF_SECRET")
	claims := RefereshTokenClaim{
		StandardClaims: jwt.StandardClaims{
			// set token lifetime in timestamp
			ExpiresAt: time.Now().Add(time.Hour * 2000).Unix(),
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
	token, err := tokenString.SignedString(refreshSecretKey)

	return token, err
}
func VerifyAccessToken(tokenString string) (AcessTokenClaim, error) {
	claims := &AcessTokenClaim{}
	_, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return accessSecretKey, nil
	})
	if err != nil {
		return *claims, err
	}
	return *claims, nil
}
func VerifyRefereshToken(tokenString string) (*RefereshTokenClaim, error) {
	claims := &RefereshTokenClaim{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return accessSecretKey, nil
	})
	if !token.Valid {
		return nil, fmt.Errorf("token is invalid")
	}
	return claims, err
}

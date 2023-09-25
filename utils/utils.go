package utils

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"golang.org/x/crypto/bcrypt"
)

func VerifyPassword(password, hashedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

func generateRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	var seededRand *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}
func IsTimeBetween(Time time.Time, timeStart int, timeEnd int) bool {
	hours, minutes, _ := Time.Clock()
	currUTCTimeInString := fmt.Sprintf("%d%02d", hours, minutes)
	currUTCTime, err := strconv.ParseInt(currUTCTimeInString, 0, 64)
	if err != nil {
		panic(err)
	}
	if int(currUTCTime) > timeStart && int(currUTCTime) < timeEnd {
		return true
	}
	return false
}

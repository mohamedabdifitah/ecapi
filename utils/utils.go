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
func IsTimeBetween(Time time.Time, timeStart float32, timeEnd float32) bool {
	hours, minutes, _ := Time.Clock()
	joinHrMinutes := fmt.Sprintf("%d.%d", hours, minutes)
	currUTCTime, err := strconv.ParseFloat(joinHrMinutes, 32)
	if err != nil {
		return false
	}
	currenTime := float32(currUTCTime)
	if currenTime > timeStart && currenTime < timeEnd {
		return true
	}
	return false

}
func GenerateIDs(length int) string {
	const charset = "0123456789"
	var seededRand *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

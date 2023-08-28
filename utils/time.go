package utils

import (
	"fmt"
	"strconv"
	"time"
)

func IsTimeBetween(timeStart int, timeEnd int) bool {
	hours, minutes, _ := time.Now().Clock()
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

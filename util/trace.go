package util

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

const charset = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

func GenerateRandomString(length int) string {
	charsetTime := fmt.Sprintf("%s%d", charset, time.Now().UnixNano()/int64(time.Millisecond))

	rand.Seed(time.Now().UnixNano())
	var sb strings.Builder
	for i := 0; i < length; i++ {
		sb.WriteByte(charsetTime[rand.Intn(len(charsetTime))])
	}
	return sb.String()
}

func GenerateTraceId(length int) string {
	currentTimeMillis := time.Now().UnixNano() / int64(time.Millisecond)
	randomString := GenerateRandomString(length)
	return randomString + fmt.Sprint(currentTimeMillis)
}

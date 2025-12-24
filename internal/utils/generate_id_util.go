package utils

import (
	"fmt"
	"math/rand"
	"time"
)

func Generate13DigitID() string {

	rand.Seed(time.Now().UnixNano())

	timestamp := time.Now().UnixMilli()

	randomPart := rand.Intn(1000)

	return fmt.Sprintf("%010d%03d", timestamp%10000000000, randomPart)
}

func GenerateOrderID() string {
	return fmt.Sprintf("TX-%s-%d", time.Now().Format("20060102"), time.Now().UnixNano())
}

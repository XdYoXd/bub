package src

import (
	"math/rand"
	"time"
)

func RandomString(strlen int) string {
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)

	const chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	result := make([]byte, strlen)
	for i := 0; i < strlen; i++ {
		result[i] = chars[r1.Intn(len(chars))]
	}

	return string(result)
}

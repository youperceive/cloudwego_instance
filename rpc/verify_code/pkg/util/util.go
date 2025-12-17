package util

import (
	"math/rand"
	"strconv"
)

func GenerateCode() string {
	return strconv.Itoa(rand.Intn(900000) + 100000)
}

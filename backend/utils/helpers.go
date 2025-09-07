package utils

import (
	"math/rand"
)

func RandomIndex(max int) int {
    return rand.Intn(max)
}

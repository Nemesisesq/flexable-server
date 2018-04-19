package utils

import (
	"math/rand"
	"time"
)

func random(min, max int) int {
	return rand.Intn(max-min) + min
}

func RandomRange(min, max int) int {
	rand.Seed(time.Now().Unix())
	myrand := random(min, max)
	return myrand
}



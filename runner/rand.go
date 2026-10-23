package main

import (
	"math/rand"
	"time"
)

var Rand *rand.Rand

func init() {
	source := rand.NewSource(time.Now().Unix())
	Rand = rand.New(source)
}

func RandBool() bool {
	return Rand.Int63() & 0x1 == 0x1
}
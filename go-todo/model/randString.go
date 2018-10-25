package model

import (
	"math/rand"
)

const (
	// randomな文字生成用
	randomLetters       = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	randomLetterIdxBits = 6
	randomLetterIdxMask = 1<<randomLetterIdxBits - 1
	randomStringLength  = 32
)

func RandString() string {
	b := make([]byte, randomStringLength)
	for i := 0; i < randomStringLength; {
		idx := int(rand.Int63() & randomLetterIdxMask)
		if idx < len(randomLetters) {
			b[i] = randomLetters[idx]
			i++
		}
	}
	return string(b)
}

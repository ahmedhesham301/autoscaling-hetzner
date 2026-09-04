package random

import "math/rand/v2"

const chars = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

// TODO use strings.Builder
func AddRandomLetters(name string) string {
	name += "-"
	for i := 0; i < 5; i++ {
		name += string(chars[rand.IntN(len(chars))])
	}
	return name
}

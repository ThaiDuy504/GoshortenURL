package shortener

import "math/rand"

const (
	ALPHABET = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	BASE = len(ALPHABET)
	ENCODE_LEN = 6	
)

func Encode() string {
	encoded := ""
	for i := 0; i < ENCODE_LEN; i++ {
		encoded = encoded + string(ALPHABET[rand.Intn(BASE)])
	}
	return encoded
}

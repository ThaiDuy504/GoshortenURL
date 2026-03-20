package shortener

import (
	b64 "encoding/base64"
	"os"
	"strconv"
)

func Encode(url string) string {
	encodeLen := 6
	if v := os.Getenv("ENCODE_LEN"); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			encodeLen = n
		}
	}
	encoded := b64.StdEncoding.EncodeToString([]byte(url))
	if encodeLen > len(encoded) {
		encodeLen = len(encoded)
	}
	return encoded[:encodeLen]
}

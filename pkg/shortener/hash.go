package shortener

import (
	b64 "encoding/base64"
)

const (
	ENCODE_LEN = 6	
)

func Encode(url string) string {
	return b64.StdEncoding.EncodeToString([]byte(url))[:ENCODE_LEN]
}

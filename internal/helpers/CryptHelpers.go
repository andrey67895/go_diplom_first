package helpers

import (
	"crypto/hmac"
	"crypto/sha256"
	"fmt"
)

func EncodeHash(value string) string {
	h := hmac.New(sha256.New, []byte("KEY123!"))
	h.Write([]byte(value))
	return fmt.Sprintf("%x", h.Sum(nil))
}

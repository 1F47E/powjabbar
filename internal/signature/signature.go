package signature

import (
	"crypto/hmac"
	"crypto/sha256"
)

func Sign(data, key, salt []byte) []byte {
	h := hmac.New(sha256.New, key)
	h.Write(data)
	h.Write(salt)
	return h.Sum(nil)
}

func Verify(data, key, salt, expectedSignature []byte) bool {
	return hmac.Equal(Sign(data, key, salt), expectedSignature)
}

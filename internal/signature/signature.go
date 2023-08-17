package signature

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"io"
)

func GenerateSignature() ([]byte, error) {
	s := make([]byte, 32)
	_, err := io.ReadFull(rand.Reader, s)
	if err != nil {
		return nil, err
	}
	return s, nil
}

func Sign(data, key, salt []byte) []byte {
	h := hmac.New(sha256.New, key)
	h.Write(data)
	h.Write(salt)
	return h.Sum(nil)
}

func Verify(data, key, salt, expectedSignature []byte) bool {
	return hmac.Equal(Sign(data, key, salt), expectedSignature)
}

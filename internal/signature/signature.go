package signature

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"io"
)

// HMAC

func GenerateSalt(length int) ([]byte, error) {
	salt := make([]byte, length)
	_, err := io.ReadFull(rand.Reader, salt)
	if err != nil {
		return nil, err
	}
	return salt, nil
}

func GenerateHMAC(data, key, salt []byte) []byte {
	h := hmac.New(sha256.New, key)
	h.Write(data)
	h.Write(salt)
	return h.Sum(nil)
}

func VerifyHMAC(data, key, salt, expectedHMAC []byte) bool {
	return hmac.Equal(GenerateHMAC(data, key, salt), expectedHMAC)
}

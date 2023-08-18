package signature

import (
	"crypto/hmac"
	"crypto/sha256"
)

type HMACSignature struct {
	key []byte
}

func New(key []byte) *HMACSignature {
	return &HMACSignature{key: key}
}

func (s *HMACSignature) Sign(data, salt []byte) []byte {
	h := hmac.New(sha256.New, s.key)
	h.Write(data)
	h.Write(salt)
	return h.Sum(nil)
}

func (s *HMACSignature) Verify(data, salt, expectedSignature []byte) bool {
	return hmac.Equal(s.Sign(data, salt), expectedSignature)
}

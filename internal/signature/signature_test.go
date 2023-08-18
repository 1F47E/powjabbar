package signature

import (
	"testing"
)

func TestHMACSignature(t *testing.T) {
	key := []byte("supersecretkey")
	signer := New(key)

	type item struct {
		data, salt    []byte
		expectedValid bool
	}
	tests := map[string]item{}
	tests["Valid signature"] = item{
		data:          []byte("data"),
		salt:          []byte("salt"),
		expectedValid: true,
	}

	validSignature := signer.Sign([]byte("data"), []byte("salt"))

	tests["Invalid data"] = item{
		data:          []byte("wrongdata"),
		salt:          []byte("salt"),
		expectedValid: false,
	}

	tests["Invalid salt"] = item{
		data:          []byte("data"),
		salt:          []byte("wrongsalt"),
		expectedValid: false,
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			if valid := signer.Verify(tt.data, tt.salt, validSignature); valid != tt.expectedValid {
				t.Errorf("Expected validity: %v, got %v", tt.expectedValid, valid)
			}
		})
	}
}

func BenchmarSign(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		key := []byte("supersecretkey")
		signer := New(key)
		data := []byte("data")
		salt := []byte("salt")

		_ = signer.Sign(data, salt)
	}
}

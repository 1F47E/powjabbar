package challenge

import (
	"bytes"
	"encoding/base64"
	"encoding/binary"
	"testing"
	"time"

	"github.com/1F47E/powjabbar/internal/signature"
	"github.com/1F47E/powjabbar/internal/utils"
	"github.com/stretchr/testify/assert"
)

func TestGenerateChallenge(t *testing.T) {
	signatureKey := []byte("signature_key")
	tests := []struct {
		name          string
		difficulty    int
		signatureKey  []byte
		expectedError error
	}{
		{
			name:          "difficulty 4",
			difficulty:    4,
			signatureKey:  signatureKey,
			expectedError: nil,
		},
		{
			name:          "difficulty 5",
			difficulty:    5,
			signatureKey:  signatureKey,
			expectedError: nil,
		},
		{
			name:          "difficulty 6",
			difficulty:    6,
			signatureKey:  signatureKey,
			expectedError: nil,
		},
		{
			name:          "zero difficulty",
			difficulty:    0,
			signatureKey:  signatureKey,
			expectedError: ErrChallengeMinDifficulty,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			challenge, err := GenerateChallenge(tt.difficulty, tt.signatureKey)

			assert.Equal(t, tt.expectedError, err)
			if err != nil {
				return
			}

			assert.Equal(t, utils.DifficultyToCriteria(tt.difficulty), challenge.Criteria)

			data, err := base64.StdEncoding.DecodeString(challenge.Data)
			assert.NoError(t, err)
			assert.Equal(t, lenDifficulty+lenTimestamp+lenNonce+lenSignature, len(data))
			assert.Equal(t, tt.difficulty, int(data[0]))

			timestamp := int64(binary.BigEndian.Uint64(data[lenDifficulty : lenDifficulty+lenTimestamp]))
			assert.True(t, time.Since(time.UnixMicro(timestamp)) <= 5*time.Second)

			nonce := data[lenDifficulty+lenTimestamp : lenDifficulty+lenTimestamp+lenNonce]
			signatureData := data[lenDifficulty : lenDifficulty+lenTimestamp+lenNonce]
			expectedSignature := signature.Sign(signatureData, tt.signatureKey, nonce)
			assert.True(t, signature.Verify(signatureData, tt.signatureKey, nonce, expectedSignature))

			assert.True(t, bytes.Equal(expectedSignature, data[lenDifficulty+lenTimestamp+lenNonce:]))
		})
	}
}

package challenge

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// simulates the client pow
func solveChallenge(data, criteria string) (string, string, string) {
	for i := 0; ; i++ {
		toSign := fmt.Sprintf("%s%d", data, i)
		hash := sha256.Sum256([]byte(toSign))
		if hex.EncodeToString(hash[:])[:len(criteria)] == criteria {
			return data, fmt.Sprintf("%d", i), hex.EncodeToString(hash[:])
		}
	}
}

func TestDynamicChallengeSolution(t *testing.T) {
	signatureKey := []byte("secret")

	testCases := []struct {
		name        string
		difficulty  int
		timelimit   time.Duration
		expectError error
	}{
		{
			name:       "Difficulty 4",
			difficulty: 4,
			timelimit:  1 * time.Second,
		},
		{
			name:       "Difficulty 5",
			difficulty: 5,
			timelimit:  1 * time.Second,
		},
		{
			name:        "timeout",
			difficulty:  4,
			timelimit:   1 * time.Microsecond,
			expectError: ErrTimelimitExceed,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			challenge, err := GenerateChallenge(tc.difficulty, signatureKey)
			assert.NoError(t, err)

			// solver simulation
			data, addedValue, hash := solveChallenge(challenge.Data, challenge.Criteria)

			sol, err := Deserialize(data)
			assert.NoError(t, err)

			valid, err := sol.Verify(data, addedValue, hash, signatureKey, tc.timelimit)
			if tc.expectError != nil {
				assert.Equal(t, tc.expectError, err)
			} else {
				assert.True(t, valid)
				assert.NoError(t, err)

			}
		})
	}
}

func TestSolutionErrors(t *testing.T) {
	signatureKey := []byte("secret")

	tests := []struct {
		name             string
		data             string
		addedValue       string
		solution         string
		timelimit        time.Duration
		timecreated      time.Time
		signatureKey     []byte
		expectedError    error
		expectValid      bool
		expectedCriteria string
	}{

		{
			name:          "Invalid difficulty",
			data:          "BgAGAyHixfjM64MRCTTjnhknYDWJqONKZ67Mx2aQEs0+UJSxfyoMlfoYMp6GsL4BaQ==",
			addedValue:    "3782172",
			solution:      "0d6035b04d0724181013ed4527a11ea082a35792a56e8d3d3e53e273000",
			signatureKey:  signatureKey,
			timelimit:     1 * time.Second,
			expectedError: ErrInvalidDifficulty,
			expectValid:   false,
		},
		{
			name:          "Invalid hash",
			data:          "BgAGAyHixfjM64MRCTTjnhknYDWJqONKZ67Mx2aQEs0+UJSxfyoMlfoYMp6GsL4BaQ==",
			addedValue:    "3782172",
			solution:      "000000d6035b04d0724181013ed4527a11ea082a35792a56e8d3d3e53e273000",
			signatureKey:  signatureKey,
			timelimit:     1 * time.Second,
			expectedError: ErrInvalidHash,
			expectValid:   false,
		},
		{
			name:             "Invalid signature",
			data:             "BAAGAyHc1PRMsIhKEpUtqVuYxBAjshlPSxD27u0aGtC4sQUL+eEXN6IssvJLBEmTmA==",
			addedValue:       "31473",
			solution:         "0000a2951eaa5fac1931aca851c04d57f1dba139c2a4a9fa68dd6ab29d0cf0ee",
			signatureKey:     []byte{0x00},
			timelimit:        1 * time.Second,
			expectedError:    ErrInvalidSignature,
			expectValid:      false,
			expectedCriteria: "0000",
		},
		{
			name:             "Invalid difficulty",
			data:             "BQAGAyHiLcs0tmSrchJtKROja9n682ODp05UY9NnrAujm284zYKne7j5t3cIA0F44g==",
			addedValue:       "968267",
			solution:         "1111155cf28f9c489642c9c08e1d2f294df86d2868a0f0a11f284a2043121d34",
			signatureKey:     signatureKey,
			timelimit:        1 * time.Second,
			expectedError:    ErrInvalidDifficulty,
			expectValid:      false,
			expectedCriteria: "00000",
		},
		{
			name:             "Timeout exceeded",
			data:             "BgAGAyHixfjM64MRCTTjnhknYDWJqONKZ67Mx2aQEs0+UJSxfyoMlfoYMp6GsL4BaQ==",
			addedValue:       "3782172",
			solution:         "000000d6035b04d0724181013ed4527a11ea082a35792a56e8d3d3e53e273876",
			signatureKey:     signatureKey,
			timelimit:        1 * time.Nanosecond, // demo case is in the past anyways, should fail
			expectedError:    ErrTimelimitExceed,
			expectValid:      false,
			expectedCriteria: "000000",
		},
		{
			name:          "Broken Base64",
			data:          "!!InvalidBase64!!",
			addedValue:    "0",
			solution:      "wronghash",
			signatureKey:  signatureKey,
			timelimit:     1 * time.Second,
			expectedError: ErrInvalidData,
			expectValid:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			solution, err := Deserialize(tt.data)
			if err != nil {
				if !errors.Is(err, tt.expectedError) {
					t.Fatalf("Unexpected error during deserialization: %v", err)
				}
				return
			}

			assert.NoError(t, err)

			isValid, err := solution.Verify(tt.data, tt.addedValue, tt.solution, tt.signatureKey, tt.timelimit)

			if tt.expectedError != nil {
				assert.Equal(t, tt.expectedError, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.expectValid, isValid)
		})
	}
}

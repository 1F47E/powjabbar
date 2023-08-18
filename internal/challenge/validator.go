// Solution validator
// The idea is to verify very cheap and fail fast
// Step 1 - verify that solution is valid via hash, this is also checks difficulty
// Step 2 - verify signature - our data was not tampered with
// Step 3 - verify timeframe

package challenge

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/1F47E/powjabbar/internal/signature"
	"github.com/1F47E/powjabbar/internal/utils"
)

type Solution struct {
	criteria   string
	timestamp  int64
	nonce      []byte
	signedData []byte
	signature  []byte
}

// deserialize data payload into the struct
func Deserialize(data string) (*Solution, error) {
	var err error

	// decode data to binary
	binData, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return nil, ErrInvalidData
	}
	if len(binData) != lenData {
		return nil, ErrInvalidDataLen
	}

	// extract data from binary
	var pos int
	difficulty := int(binData[pos])
	pos++

	timestamp := int64(binary.BigEndian.Uint64(binData[pos : pos+lenTimestamp]))
	pos += lenTimestamp

	nonce := binData[pos : pos+lenNonce]

	signedData := make([]byte, lenSignedData)
	copy(signedData, binData[lenDifficulty:lenDifficulty+lenSignedData])
	signature := binData[len(binData)-lenSignature:]

	return &Solution{
		criteria:   utils.DifficultyToCriteria(difficulty),
		timestamp:  timestamp,
		nonce:      nonce,
		signedData: signedData,
		signature:  signature,
	}, nil
}

// verify that submitted solution is valid
func (s *Solution) Verify(data, addedvalue, solution string, signatureKey []byte, timelimit time.Duration) (bool, error) {
	// difficulty
	criteriaMet := solution[:len(s.criteria)] == s.criteria
	if !criteriaMet {
		return false, ErrInvalidDifficulty
	}

	// hash
	toHash := fmt.Sprintf("%s%s", data, addedvalue)
	hash := sha256.Sum256([]byte(toHash))
	if hex.EncodeToString(hash[:]) != solution {
		return false, ErrInvalidHash
	}

	// signature
	signer := signature.New(signatureKey)
	if !signer.Verify(s.signedData, s.nonce, s.signature) {
		return false, ErrInvalidSignature
	}

	// timelimit
	timestamp := time.UnixMicro(s.timestamp)
	if time.Since(timestamp) > timelimit {
		return false, ErrTimelimitExceed
	}
	return true, nil
}

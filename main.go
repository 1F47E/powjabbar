package powjabbar

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"io"
	"time"

	"github.com/1F47E/powjabbar/internal/signature"
	"github.com/1F47E/powjabbar/internal/validator"
)

const (
	lenDifficulty = 1
	lenTimestamp  = 8
	lenNonce      = 8
	lenSignature  = 32
)

// signatureKey - secret key used to sign the data payload.
// Recommended to be at least 32 bytes long
type PowJabbar struct {
	signatureKey []byte
}

func NewPowJabbar(signatureKey []byte) *PowJabbar {
	return &PowJabbar{
		signatureKey: signatureKey,
	}
}

// GenerateSignature generates a random signature key 32 bytes long
// This key will be used to sign the data payload
// Key is generated using crypto/rand and may return error sometimes
func GenerateSignature() ([]byte, error) {
	return signature.GenerateSignature()
}

// Challenge data format
// DIFFICULTY|TIMESTAMP|NONCE|SIGNATURE
// 4|1692065996206899|1692065996206899|7814f500270011d762ad116acd45c97a455e079a9d958746cb8e813a7828ed81
// 1 byte | 8 bytes| 8 byte | 32 bytes
// total 49 bytes
type Challenge struct {
	Data     string
	Criteria string
}

// Generate challenge for the client
// difficulty - number of leading zeroes in the hash
// The more leading zeros in a hash, the more difficult it is to find the solution
// Choosen difficulty level will be baked into the data payload
func (p *PowJabbar) GenerateChallenge(difficulty int) (*Challenge, error) {
	// TODO: change to const type difficulty
	if difficulty == 0 {
		return nil, fmt.Errorf("difficulty must be greater than 0")
	}

	criteria := make([]rune, difficulty)
	for i := 0; i < difficulty; i++ {
		criteria[i] = '0'
	}

	nonce := make([]byte, lenNonce)
	_, err := io.ReadFull(rand.Reader, nonce)
	if err != nil {
		return nil, fmt.Errorf("generate nonce failed: %v", err)
	}

	// convert timestamp to binary
	timestamp := time.Now().UnixMicro()
	bTimestamp := make([]byte, lenTimestamp)
	binary.BigEndian.PutUint64(bTimestamp, uint64(timestamp))

	// sign data
	signData := make([]byte, lenTimestamp+lenNonce)
	copy(signData, bTimestamp)
	copy(signData[lenTimestamp:], nonce)
	bSignature := signature.Sign(signData, p.signatureKey, nonce)

	// assemble result data
	data := make([]byte, lenDifficulty+len(signData)+lenSignature)
	pos := 0
	data[pos] = byte(difficulty)
	pos = lenDifficulty
	copy(data[pos:], signData)
	pos += len(signData)
	copy(data[pos:], bSignature)

	dataStr := base64.StdEncoding.EncodeToString(data)

	return &Challenge{
		Data:     dataStr,
		Criteria: string(criteria),
	}, nil
}

// Verify clients solution
// data - data payload given to the client
// value - added value added by the client
// hash - hash of the data + added value
func (p *PowJabbar) VerifySolution(data, value, hash string, timelimit time.Duration) (bool, error) {
	// TODO: add reason for failed validations
	sol, err := validator.Deserialize(data)
	if err != nil {
		return false, fmt.Errorf("solution data is broken: %v", err)
	}

	// Step by steps check, fail early and fast
	if !sol.VerifySolution(data, value, hash) {
		return false, nil
	}

	// check time limit
	// TODO embed timelimit into the data payload
	if !sol.VerifyTimelimit(timelimit) {
		return false, nil
	}

	// verify signature
	if !sol.VerifySignature(p.signatureKey) {
		return false, nil
	}
	return true, nil
}

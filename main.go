package powjabbar

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"io"
	"time"

	"github.com/1F47E/go-pow-jabbar/internal/signature"
	"github.com/1F47E/go-pow-jabbar/internal/solution"
)

const (
	lenDifficulty = 1
	lenTimestamp  = 8
	lenNonce      = 8
	lenSignature  = 32
)

// Challenge data format
// DIFFICULTY|TIMESTAMP|NONCE|SIGNATURE
// 4|1692065996206899|1692065996206899|7814f500270011d762ad116acd45c97a455e079a9d958746cb8e813a7828ed81
// 1 byte | 8 bytes| 8 byte | 32 bytes
type Challenge struct {
	Data     string
	Criteria string
}

// Generate challenge for the client
// difficulty - number of leading zeroes in the hash
// The more leading zeros in a hash, the more difficult it is to find the solution
// Choosen difficulty level will be baked into the data payload
func GenerateChallenge(difficulty int, key []byte) (*Challenge, error) {
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
	bSignature := signature.GenerateHMAC(signData, key, nonce)

	// assemble result data
	// p for readability
	data := make([]byte, lenDifficulty+len(signData)+lenSignature)
	p := 0
	data[p] = byte(difficulty)
	p = lenDifficulty
	copy(data[p:], signData)
	p += len(signData)
	copy(data[p:], bSignature)

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
func VerifySolution(data, value, hash string, signatureKey []byte, timelimit time.Duration) (bool, error) {
	sol, err := solution.Deserialize(data)
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
	if !sol.VerifySignature(signatureKey) {
		return false, nil
	}
	return true, nil
}

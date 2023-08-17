package challenge

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"io"
	"time"

	"github.com/1F47E/powjabbar/internal/signature"
	"github.com/1F47E/powjabbar/internal/utils"
)

// Challenge data format
// DIFFICULTY|TIMESTAMP|NONCE|SIGNATURE
// 4|1692065996206899|1692065996206899|7814f500270011d762ad116acd45c97a455e079a9d958746cb8e813a7828ed81
// 1 byte | 8 bytes| 8 byte | 32 bytes
// total 49 bytes
type Challenge struct {
	Data     string
	Criteria string
}

func GenerateChallenge(difficulty int, signatureKey []byte) (*Challenge, error) {
	if difficulty == 0 {
		return nil, fmt.Errorf("difficulty must be greater than 0")
	}

	// generate random nonce
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
	bSignature := signature.Sign(signData, signatureKey, nonce)

	// assemble result data payload
	data := make([]byte, lenDifficulty+len(signData)+lenSignature)
	pos := 0
	data[pos] = byte(difficulty)
	pos = lenDifficulty
	copy(data[pos:], signData)
	pos += len(signData)
	copy(data[pos:], bSignature)
	dataEncoded := base64.StdEncoding.EncodeToString(data)

	return &Challenge{
		Data:     dataEncoded,
		Criteria: utils.DifficultyToCriteria(difficulty),
	}, nil
}

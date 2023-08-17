// Solution validator
// The idea is to verify very cheap and fail fast
// Step 1 - verify that solution is valid via hash, this is also checks difficulty
// Step 2 - verify signature - our data was not tampered with
// Step 3 - verify timeframe

package validator

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/1F47E/powjabbar/internal/signature"
)

type Validator struct {
	Criteria   string
	Timestamp  int64
	Nonce      []byte
	SignedData []byte
	Signature  []byte
}

// deserialize data payload into the struct
func Deserialize(data string) (*Validator, error) {
	// TODO: check data len and fail early
	var err error
	bindata, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return nil, err
	}
	const (
		lenDifficulty = 1
		lenTimestamp  = 8
		lenNonce      = 8
		lenSignedData = lenTimestamp + lenNonce
		lenSignature  = 32
	)

	s := Validator{}

	// get difficulty from the data first byte
	// and create criteria string, witch is a string of zeros
	pos := 0
	difficulty := int(bindata[pos])
	pos++
	criteria := make([]rune, difficulty)
	for i := 0; i < difficulty; i++ {
		criteria[i] = '0'
	}
	s.Criteria = string(criteria)

	s.Timestamp = int64(binary.BigEndian.Uint64(bindata[pos : pos+lenTimestamp]))
	pos += lenTimestamp

	s.Nonce = bindata[pos : pos+lenNonce]

	// collect data to sign (timestamp + nonce)
	signedData := make([]byte, lenSignedData)
	copy(signedData, bindata[lenDifficulty:lenDifficulty+lenSignedData])
	s.SignedData = signedData

	s.Signature = bindata[len(bindata)-lenSignature:]
	return &s, nil
}

// verify that submitted solution is valid
func (v *Validator) VerifySolution(data, addedvalue, solution string) bool {
	// verify that difficulty criteria is met
	criteriaMet := solution[:len(v.Criteria)] == v.Criteria
	if !criteriaMet {
		fmt.Println("criteria not met")
		return false
	}
	toHash := fmt.Sprintf("%s%s", data, addedvalue)

	hash := sha256.Sum256([]byte(toHash))
	return hex.EncodeToString(hash[:]) == solution
}

// verify integrity of the solution
func (s *Validator) VerifySignature(key []byte) bool {
	return signature.Verify(s.SignedData, key, s.Nonce, s.Signature)
}

// check if it was solved within the timelimit
// TODO: bake timelimit to the date
func (s *Validator) VerifyTimelimit(timelimit time.Duration) bool {
	timestamp := time.UnixMicro(s.Timestamp)
	return time.Since(timestamp) < timelimit
}

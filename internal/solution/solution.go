// Solution verifier
// The idea is to verify very cheap and fail fast
// Step 1 - verify that solution is valid via hash
// Step 2 - verify that valid solution is not tempered with (signature)
// Step 3 - verify that solution is not too old (timestamp)

package solution

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/1F47E/go-pow-jabbar/internal/signature"
)

// TODO: rename to validator
type Solution struct {
	// data string

	// fields extracted from data
	Criteria   string
	Timestamp  int64
	Nonce      []byte
	SignedData []byte
	Signature  []byte
}

// verify that submitted solution is valid
func (s *Solution) VerifySolution(data, addedvalue, solution string) bool {
	// verify that difficulty criteria is met
	criteriaMet := solution[:len(s.Criteria)] == s.Criteria
	if !criteriaMet {
		fmt.Println("criteria not met")
		return false
	}
	toHash := fmt.Sprintf("%s%s", data, addedvalue)

	hash := sha256.Sum256([]byte(toHash))
	return hex.EncodeToString(hash[:]) == solution
}

// deserialize data payload into the struct
func Deserialize(data string) (*Solution, error) {
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

	s := Solution{}

	// get difficulty from the data first byte
	p := 0
	difficulty := int(bindata[p])
	p++
	criteria := make([]rune, difficulty)
	for i := 0; i < difficulty; i++ {
		criteria[i] = '0'
	}
	s.Criteria = string(criteria)

	s.Timestamp = int64(binary.BigEndian.Uint64(bindata[p : p+lenTimestamp]))
	p += lenTimestamp

	// get nonce
	s.Nonce = bindata[p : p+lenNonce]

	// get data to sign (timestamp + nonce)
	signedData := make([]byte, lenSignedData)
	copy(signedData, bindata[lenDifficulty:lenDifficulty+lenSignedData])
	s.SignedData = signedData

	// get signature
	s.Signature = bindata[len(bindata)-lenSignature:]
	return &s, nil
}

// verify integrity of the solution
func (s *Solution) VerifySignature(key []byte) bool {
	return signature.VerifyHMAC(s.SignedData, key, s.Nonce, s.Signature)
}

// If solution is valid, check if it was solved within the timelimit
// TODO: bake timelimit to the date
func (s *Solution) VerifyTimelimit(timelimit time.Duration) bool {
	timestamp := time.UnixMicro(s.Timestamp)
	return time.Since(timestamp) < timelimit
}

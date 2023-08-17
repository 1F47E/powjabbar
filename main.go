package powjabbar

import (
	"fmt"
	"time"

	"github.com/1F47E/powjabbar/internal/challenge"
	// "github.com/1F47E/powjabbar/internal/validator"
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

// Generate easy challenge
// With difficulty level 4
// Estimated time to solve ~40ms
func (p *PowJabbar) GenerateChallangeEasy() (*challenge.Challenge, error) {
	return p.GenerateChallenge(4)
}

// Generate medium challenge
// With difficulty level 5
// Estimated time to solve ~100ms
func (p *PowJabbar) GenerateChallangeMedium() (*challenge.Challenge, error) {
	return p.GenerateChallenge(5)
}

// Generate hard challenge
// With difficulty level 6
// Estimated time to solve ~10+ sec
func (p *PowJabbar) GenerateChallangeHard() (*challenge.Challenge, error) {
	return p.GenerateChallenge(6)
}

// Generate challenge for the client
// difficulty - number of leading zeroes in the hash
// The more leading zeros in a hash, the more difficult it is to find the solution
// Choosen difficulty level will be baked into the data payload
// Depending on the hardware, the time to solve the challenge may vary
// 4 - 20-50ms
// 5 - 100-200ms
// 6 - 15+ sec
func (p *PowJabbar) GenerateChallenge(difficulty int) (*challenge.Challenge, error) {
	// TODO: change to const type difficulty
	if difficulty == 0 {
		return nil, fmt.Errorf("difficulty must be greater than 0")
	}
	return challenge.GenerateChallenge(difficulty, p.signatureKey)
}

// Verify clients solution
// data - data payload given to the client
// value - added value added by the client
// hash - hash of the data + added value
func (p *PowJabbar) VerifySolution(data, value, hash string, timelimit time.Duration) (bool, error) {
	// TODO: add reason for failed validations
	sol, err := challenge.Deserialize(data)
	if err != nil {
		return false, fmt.Errorf("solution data is broken: %v", err)
	}

	// Step by steps check, fail early and fast
	return sol.Verify(data, value, hash, p.signatureKey, timelimit)
}

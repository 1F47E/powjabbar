// create a random nonce, send it to the client
// client generates

package main

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	timelimit         = 1 * time.Second
	defaultDifficulty = 4
	colorReset        = "\033[0m"
	colorRed          = "\033[31m"
	colorGreen        = "\033[32m"
	colorYellow       = "\033[33m"
	colorPurple       = "\033[35m"
	colorCyan         = "\033[36m"
)

func printSuccess(msg string) {
	fmt.Println(colorGreen, msg, colorReset)
}

func printError(msg string) {
	fmt.Println(colorRed, msg, colorReset)
}

func printInfo(msg string) {
	fmt.Println(colorCyan, msg, colorReset)
}

// HMAC

var signatureKey = []byte("secret")

func GenerateSalt(length int) ([]byte, error) {
	salt := make([]byte, length)
	_, err := io.ReadFull(rand.Reader, salt)
	if err != nil {
		return nil, err
	}
	return salt, nil
}

func GenerateHMAC(data, key, salt []byte) string {
	h := hmac.New(sha256.New, key)
	h.Write(data)
	h.Write(salt) // append salt to the data before generating the HMAC
	signature := hex.EncodeToString(h.Sum(nil))
	return signature
}

func VerifyHMAC(data, key, salt []byte, expectedHMAC string) bool {
	gen := GenerateHMAC(data, key, salt)
	return hmac.Equal([]byte(gen), []byte(expectedHMAC))
}

// POW CHALLENGE
// Server's challenge
type challenge struct {
	data      string
	nonce     string // to use as a salt
	criteria  string
	signature string
}

func NewChallenge(difficulty int, key []byte) (*challenge, error) {
	if difficulty == 0 {
		return nil, fmt.Errorf("difficulty must be greater than 0")
	}
	nonce := make([]byte, 1)
	_, err := io.ReadFull(rand.Reader, nonce)
	if err != nil {
		return nil, fmt.Errorf("generate nonce failed: %v", err)
	}
	nonceStr := fmt.Sprintf("%d", nonce[0])
	timestamp := fmt.Sprintf("%d", time.Now().UnixMicro())
	toSign := fmt.Sprintf("%s|%s", timestamp, nonceStr)
	signature := GenerateHMAC([]byte(toSign), key, []byte(nonceStr))

	// final data looks like
	// TIMESTAMP|NONCE|SIGNATURE
	data := fmt.Sprintf("%s|%s", toSign, signature)

	// gen difficulty
	symbol := "0"
	criteria := ""
	for i := 0; i < difficulty; i++ {
		criteria += symbol
	}

	return &challenge{
		nonce:     nonceStr,
		data:      data,
		criteria:  criteria, //  hash must start with this string
		signature: signature,
	}, nil
}

// Client's solution
type solution struct {
	data       string
	addedValue string
	hash       string

	// fields extracted from data
	timestamp int64
	nonce     string
	signature string
}

func (s *solution) deserialize() error {
	parts := strings.Split(s.data, "|")
	if len(parts) != 3 {
		return fmt.Errorf("invalid base data: %s", s.data)
	}

	ts := parts[0]
	// just to be sure the timestamp is int
	timestamp, err := strconv.ParseInt(ts, 10, 64)
	if err != nil {
		return fmt.Errorf("invalid timestamp, not int: %v", err)
	}
	s.timestamp = timestamp
	s.nonce = parts[1]
	// just to be sure the nonce is uint
	_, err = strconv.ParseUint(s.nonce, 10, 64)
	if err != nil {
		return fmt.Errorf("invalid nonce, not int: %v", err)
	}
	s.signature = parts[2]
	return nil
}

func (s *solution) verifySolution() bool {
	data := fmt.Sprintf("%s%s", s.data, s.addedValue)
	hash := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hash[:]) == s.hash
}

func (s *solution) verifySignature(key []byte) bool {
	dataToSign := fmt.Sprintf("%d|%s", s.timestamp, s.nonce)
	return VerifyHMAC([]byte(dataToSign), key, []byte(s.nonce), s.signature)
}

func (s *solution) verifyTimelimit() bool {
	timestamp := time.UnixMicro(s.timestamp)

	// get diff just for debug logs
	diff := time.Since(timestamp)
	diffMs := diff.Milliseconds()
	if diff > timelimit {
		printError(fmt.Sprintf("timestamp is too old, took %v/%v", diff, timelimit))
		return false
	}
	printInfo(fmt.Sprintf("elapsed %v/%v, diff: %v\n", time.Since(timestamp), timelimit, diffMs))

	return true
}

// client side POV
func solveChallenge(data, criteria string) solution {
	for i := 0; ; i++ {
		dataPlus := fmt.Sprintf("%s%d", data, i)

		hash := sha256.Sum256([]byte(dataPlus))
		if hex.EncodeToString(hash[:])[:len(criteria)] == criteria {
			return solution{
				data:       data,
				addedValue: fmt.Sprintf("%d", i),
				hash:       hex.EncodeToString(hash[:]),
			}
		}
	}
}

func main() {
	// try get difficulty from flag
	args := os.Args[1:]
	difficulty := defaultDifficulty
	var err error
	if len(args) != 1 {
		fmt.Printf("using default difficulty: %d\n", difficulty)
	} else {
		difficulty, err = strconv.Atoi(args[0])
		if err != nil {
			log.Fatalf("invalid difficulty: %v. expecting integer", err)
		}
		fmt.Printf("using difficulty: %s\n", args[0])
	}

	// ===== Server generates the challenge

	now := time.Now()
	c, err := NewChallenge(difficulty, signatureKey)
	if err != nil {
		log.Fatalf("generate challenge failed: %v", err)
	}
	fmt.Printf("Challenge Generated: %+v, took: %v\n", c, time.Since(now))

	// ===== Client solves the challenge

	now = time.Now()
	sol := solveChallenge(c.data, c.criteria)
	err = sol.deserialize()
	if err != nil {
		log.Fatalf("solution deserialization error: %v", err)
	}
	msg := fmt.Sprintf("Solution found:\n%+v\ntook: %v\n", sol, time.Since(now))
	printInfo(msg)

	// ===== Server verifies the solution

	now = time.Now()

	// check time limit
	// if timelimit fails - do not bother with other checks
	if sol.verifyTimelimit() {
		printSuccess("Request is within time limit!")
	} else {
		printError("Request is out of time limit!")
	}

	// check hash
	if sol.verifySolution() {
		printSuccess("Solution is valid!")
	} else {
		printError("Solution is invalid!")
	}

	// verify signature
	if sol.verifySignature(signatureKey) {
		printSuccess("Signature is valid!")
	} else {
		printError("Signature is invalid!")
	}

	printInfo(fmt.Sprintf("Verification, took: %v\n", time.Since(now)))
}

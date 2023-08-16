package main

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"time"
)

var signatureKey = []byte("secret")

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

func GenerateSalt(length int) ([]byte, error) {
	salt := make([]byte, length)
	_, err := io.ReadFull(rand.Reader, salt)
	if err != nil {
		return nil, err
	}
	return salt, nil
}

func GenerateHMAC(data, key, salt []byte) []byte {
	h := hmac.New(sha256.New, key)
	h.Write(data)
	h.Write(salt)
	return h.Sum(nil)
}

func VerifyHMAC(data, key, salt, expectedHMAC []byte) bool {
	return hmac.Equal(GenerateHMAC(data, key, salt), expectedHMAC)
}

// POW CHALLENGE

// Server's challenge
type challenge struct {
	data     string
	criteria string
}

// data format
// TIMESTAMP|NONCE|SIGNATURE
// 1692065996206899|1692065996206899|7814f500270011d762ad116acd45c97a455e079a9d958746cb8e813a7828ed81
// 8 bytes| 8 byte | 64 bytes
func NewChallenge(difficulty int, key []byte) (*challenge, error) {
	if difficulty == 0 {
		return nil, fmt.Errorf("difficulty must be greater than 0")
	}

	nonce := make([]byte, 8)
	_, err := io.ReadFull(rand.Reader, nonce)
	if err != nil {
		return nil, fmt.Errorf("generate nonce failed: %v", err)
	}

	// convert timestamp to binary
	timestamp := time.Now().UnixMicro()
	bTimestamp := make([]byte, 8)
	binary.BigEndian.PutUint64(bTimestamp, uint64(timestamp))

	// sign data
	signData := make([]byte, 8+8) // timestamp + nonce
	copy(signData, bTimestamp)
	copy(signData[8:], nonce)
	bSignature := GenerateHMAC(signData, key, nonce)

	// assemble result data
	data := make([]byte, 16+32)
	copy(data, signData)
	copy(data[16:], bSignature)

	// gen difficulty
	runes := make([]rune, difficulty)
	for i := 0; i < difficulty; i++ {
		runes[i] = '0'
	}
	dataStr := base64.StdEncoding.EncodeToString(data)

	return &challenge{
		data:     dataStr,
		criteria: string(runes),
	}, nil
}

// Client's solution
type solution struct {
	data       string
	addedValue string
	hash       string

	// fields extracted from data
	timestamp  int64
	nonce      []byte
	signedData []byte
	signature  []byte
}

// deserialize with base64
func (s *solution) deserialize() error {
	var err error
	bindata, err := base64.StdEncoding.DecodeString(s.data)
	if err != nil {
		return err
	}
	// get timestamp from the data first 8 bytes
	s.timestamp = int64(binary.BigEndian.Uint64(bindata[0:8]))

	// get nonce
	s.nonce = bindata[8:16]

	// get data to sign (timestamp + nonce)
	signedData := make([]byte, 16)
	copy(signedData, bindata[0:16])
	s.signedData = signedData

	// get signature
	s.signature = bindata[16:]
	return nil
}

func (s *solution) verifySolution() bool {
	data := fmt.Sprintf("%s%s", s.data, s.addedValue)
	hash := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hash[:]) == s.hash
}

func (s *solution) verifySignature(key []byte) bool {
	return VerifyHMAC(s.signedData, key, s.nonce, s.signature)
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
	testB64(difficulty, signatureKey)

}

func testB64(difficulty int, signatureKey []byte) {
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
	msg := fmt.Sprintf("Solution found, took: %v\n", time.Since(now))
	msg += fmt.Sprintf("data: %+v\n", sol.data)
	msg += fmt.Sprintf("added value: %+v\n", sol.addedValue)
	msg += fmt.Sprintf("hash: %+v\n", sol.hash)
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

package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	pj "github.com/1F47E/go-pow-jabbar"
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

// client side POV
func solveChallenge(data, criteria string) (string, string, string) {
	for i := 0; ; i++ {
		toSign := fmt.Sprintf("%s%d", data, i)
		hash := sha256.Sum256([]byte(toSign))
		if hex.EncodeToString(hash[:])[:len(criteria)] == criteria {
			return data, fmt.Sprintf("%d", i), hex.EncodeToString(hash[:])
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

	now := time.Now()
	pow := pj.NewPowJabbar(signatureKey)

	// difficulty can be changed on the fly
	c, err := pow.GenerateChallenge(difficulty)
	if err != nil {
		log.Fatalf("generate challenge failed: %v", err)
	}
	fmt.Printf("Challenge Generated: %+v, took: %v\n", c, time.Since(now))

	// ===== Client solves the challenge

	now = time.Now()
	// EMULATE data coming from the client side
	solData, solValue, solHash := solveChallenge(c.Data, c.Criteria)
	msg := fmt.Sprintf("Solution found, took: %v\n", time.Since(now))
	msg += fmt.Sprintf("data: %+v\n", solData)
	msg += fmt.Sprintf("added value: %+v\n", solValue)
	msg += fmt.Sprintf("hash: %+v\n", solHash)
	printInfo(msg)

	// ===== VERIFY SOLUTION

	now = time.Now()
	success, err := pow.VerifySolution(solData, solValue, solHash, timelimit)
	if err != nil {
		log.Fatalf("solution error: %v", err)
	}
	if !success {
		log.Fatalf("solution error: %v", err)
	}
	printSuccess("Solution is valid!")

	printInfo(fmt.Sprintf("Verification, took: %v\n", time.Since(now)))

}

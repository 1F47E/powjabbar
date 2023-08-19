package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/1F47E/powjabbar"
)

const (
	defaultDifficulty = 4
	timelimit         = 2 * time.Second
)

// Create a new powjabbar instance with a secret signature key (32 bytes)
var (
	signatureKey = []byte("your-secret-signature-key")
	pow          = powjabbar.NewPowJabbar(signatureKey)
)

type ChallengeResponse struct {
	Data        string `json:"data"`
	Criteria    string `json:"criteria"`
	TimelimitMs int64  `json:"timelimit_ms"` // in ms, optional. Just so solver can fail if hit timeout
}

type SolutionRequest struct {
	Data  string `json:"data"`
	Value string `json:"value"`
	Hash  string `json:"hash"`
}

type SolutionResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error,omitempty"`
}

func handlerGetChallenge(w http.ResponseWriter, r *http.Request) {
	// get difficulty from query string
	q := r.URL.Query()
	difficulty := 4
	if d := q.Get("difficulty"); d != "" {
		fmt.Printf("difficulty: %s\n", d)
		df, err := strconv.Atoi(d)
		if err != nil {
			http.Error(w, "Invalid difficulty", http.StatusBadRequest)
			return
		}
		difficulty = df
	}
	c, err := pow.GenerateChallenge(difficulty)
	if err != nil {
		http.Error(w, "generate challenge failed", http.StatusInternalServerError)
		return
	}

	resp := ChallengeResponse{
		Data:        c.Data,
		Criteria:    c.Criteria,
		TimelimitMs: timelimit.Milliseconds(),
	}
	log.Printf("challenge generated: %+v", resp)

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(resp)
}

func handlerSolution(w http.ResponseWriter, r *http.Request) {
	// parse request data
	req := SolutionRequest{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	success, err := pow.VerifySolution(req.Data, req.Value, req.Hash, timelimit)
	resp := SolutionResponse{
		Success: success,
	}
	if err != nil {
		log.Printf("solution error: %v", err)
		resp.Error = err.Error()
	}
	log.Printf("solution success: %+v", success)

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(resp)
}

func main() {
	http.HandleFunc("/challenge", handlerGetChallenge)
	http.HandleFunc("/solution", handlerSolution)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))

	hostport := "localhost:8080"
	log.Println("Server started at", hostport)
	log.Fatal(http.ListenAndServe(hostport, nil))
}

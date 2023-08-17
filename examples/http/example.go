package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	pj "github.com/1F47E/powjabbar"
)

var (
	signatureKey = []byte("your-secret-signature-key")
	pow          = pj.NewPowJabbar(signatureKey)
	difficulty   = 4
	timelimit    = 1 * time.Second
)

type ChallengeResponse struct {
	Data     string `json:"data"`
	Criteria string `json:"criteria"`
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
	c, err := pow.GenerateChallenge(difficulty)
	if err != nil {
		http.Error(w, "generate challenge failed", http.StatusInternalServerError)
		return
	}

	resp := ChallengeResponse{
		Data:     c.Data,
		Criteria: c.Criteria,
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

	log.Println("Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

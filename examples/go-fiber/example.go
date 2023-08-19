package main

import (
	"log"
	"strconv"
	"time"

	"github.com/1F47E/powjabbar"
	"github.com/gofiber/fiber/v2"
)

const (
	defaultDifficulty = 4
	timelimit         = 2 * time.Second // for the client to solve the challenge
)

var (
	signatureKey = []byte("secret-signature-key")
	pow          = powjabbar.NewPowJabbar(signatureKey)
)

type ChallengeResponse struct {
	Data        string `json:"data"`
	Criteria    string `json:"criteria"`
	TimelimitMs int64  `json:"timelimit_ms"`
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

func handlerGetChallenge(c *fiber.Ctx) error {
	difficulty, _ := strconv.Atoi(c.Query("difficulty", "4"))

	criteria, err := pow.GenerateChallenge(difficulty)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "generate challenge failed",
		})
	}

	resp := ChallengeResponse{
		Data:        criteria.Data,
		Criteria:    criteria.Criteria,
		TimelimitMs: timelimit.Milliseconds(),
	}

	return c.JSON(resp)
}

func handlerSolution(c *fiber.Ctx) error {
	req := new(SolutionRequest)
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	success, err := pow.VerifySolution(req.Data, req.Value, req.Hash, timelimit)
	resp := SolutionResponse{
		Success: success,
	}
	if err != nil {
		resp.Error = err.Error()
	}

	return c.JSON(resp)
}

func main() {
	app := fiber.New()

	app.Get("/challenge", handlerGetChallenge)
	app.Post("/solution", handlerSolution)
	app.Static("/static", "./static")

	hostport := "localhost:8080"
	log.Println("Server started at", hostport)
	log.Fatal(app.Listen(hostport))
}

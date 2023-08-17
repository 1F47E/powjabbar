package main

import (
	"testing"
	"time"

	pj "github.com/1F47E/go-pow-jabbar"
	"github.com/stretchr/testify/assert"
)

func BenchmarkExtraction(b *testing.B) {

	// fake solution
	data := "AQAGAxUWHmkmP4VtalxZ71xaDiK+xy5wE/lqrgLg0TPzF7ey/rTAq1HEVCPG/EAMKg=="
	value := "7"
	hash := "0563442751d2b2bfb418d1d34e0db20c979961d3d3bf89dd206754b04c247a52"

	// benchmark
	b.ResetTimer()
	pow := pj.NewPowJabbar([]byte("secret"))
	for i := 0; i < b.N; i++ {
		valid, err := pow.VerifySolution(data, value, hash, time.Second)
		assert.NoError(b, err)
		assert.True(b, valid)
	}
}

func BenchmarkNewChallenge(b *testing.B) {
	pow := pj.NewPowJabbar([]byte("test"))
	for i := 0; i < b.N; i++ {
		_, err := pow.GenerateChallenge(3)
		assert.NoError(b, err)
	}
}

func TestGetChallenge(t *testing.T) {
	pow := pj.NewPowJabbar([]byte("test"))
	t.Run("Generate challenge", func(t *testing.T) {
		_, err := pow.GenerateChallenge(3)
		assert.NoError(t, err)
	})

}

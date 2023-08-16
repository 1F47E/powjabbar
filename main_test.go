package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func BenchmarkExtraction(b *testing.B) {

	// fake solution
	sol := solution{
		data:       "AAYC/hQ8DX7ISGXWiYsuuBY0LcqvV4MbPMXU/cab5DyjOdpI9Hkro90=",
		addedValue: "44927",
		hash:       "0000caae3e3fd8d25b64e4eaa982b65921f7724421776e7646f86736dcf75dd7",
	}

	// benchmark
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err := sol.deserialize()
		assert.NoError(b, err)
	}
}

func BenchmarkNewChallenge(b *testing.B) {
	var test_key = []byte("test_key")
	for i := 0; i < b.N; i++ {
		_, err := NewChallenge(4, test_key)
		assert.NoError(b, err)
	}
}

func TestGetChallenge(t *testing.T) {
	var test_key = []byte("test_key")
	// test challenge generation
	t.Run("Generate challenge", func(t *testing.T) {
		_, err := NewChallenge(3, test_key)
		assert.NoError(t, err)
	})

}

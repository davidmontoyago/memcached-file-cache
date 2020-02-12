package chunker

import (
	"math/rand"
	"time"
)

// minChunkSize is the min chunk size
const minChunkSize = 96

// maxChunkSize is the max chunk size
const maxChunkSize = 1024

// chunkSizer generates random chunk sizes between minChunkSize and maxChunkSize
type chunkSizer struct {
	randSource *rand.Rand
}

func newChunkSizer() *chunkSizer {
	randomSource := rand.NewSource(time.Now().UnixNano())
	return &chunkSizer{
		randSource: rand.New(randomSource),
	}
}

// returns a random chunk size
func (c *chunkSizer) new() int {
	chunkSize := c.randSource.Intn(maxChunkSize - minChunkSize + 1)
	return minChunkSize + chunkSize
}

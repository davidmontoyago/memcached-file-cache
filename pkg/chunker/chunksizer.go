package chunker

import (
	"math/rand"
	"time"

	xrand "golang.org/x/exp/rand"
	"gonum.org/v1/gonum/spatial/r1"
	"gonum.org/v1/gonum/stat/distmv"
)

// minChunkSize is the min chunk size
const minChunkSize = 96

// maxChunkSize is the max chunk size
const maxChunkSize = 1048576

// ChunkSizer represents an strategy for determining chunk sizes
type ChunkSizer interface {
	NextChunkSize() int
}

// --------------------------- UNIFORM distribution chunk sizer ---------------------------
type uniformChunkSizer struct {
	dist       *distmv.Uniform
	randSource xrand.Source
}

func newUniformChunkSizer() ChunkSizer {
	randomSource := xrand.NewSource(uint64(time.Now().UnixNano()))
	dist := distmv.NewUniform([]r1.Interval{{Min: minChunkSize, Max: maxChunkSize}}, randomSource)
	return &uniformChunkSizer{
		dist:       dist,
		randSource: randomSource,
	}
}

func (c *uniformChunkSizer) NextChunkSize() int {
	return int(c.dist.Rand(nil)[0])
}

// --------------------------- RANDOM chunk sizer ---------------------------
// randomChunkSizer generates random chunk sizes between minChunkSize and maxChunkSize
type randomChunkSizer struct {
	randSource *rand.Rand
}

func newRandomChunkSizer() ChunkSizer {
	randomSource := rand.NewSource(time.Now().UnixNano())
	return &randomChunkSizer{
		randSource: rand.New(randomSource),
	}
}

// returns a random chunk size
func (c *randomChunkSizer) NextChunkSize() int {
	chunkSize := c.randSource.Intn(maxChunkSize - minChunkSize + 1)
	return minChunkSize + chunkSize
}

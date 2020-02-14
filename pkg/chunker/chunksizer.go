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

// --------------------------- SKEWED distribution chunk sizer ---------------------------
type skewedChunkSizer struct {
	randSource xrand.Source
	random     *xrand.Rand
}

func newSkewedChunkSizer() ChunkSizer {
	randomSource := xrand.NewSource(uint64(time.Now().UnixNano()))
	return &skewedChunkSizer{
		randSource: randomSource,
		random:     xrand.New(randomSource),
	}
}

func (c *skewedChunkSizer) NextChunkSize() int {
	// 80% of the time, fall in the lower range 96 to 246KB
	// 20% of the time, fall in the higher range > 246KB
	frequenceRatio := (float64(20) / float64(80))
	skew := c.random.Float64()

	var chunkSize int
	if skew > frequenceRatio {
		chunkSize = randomInRange(minChunkSize, 251904-1, c.random)
	} else {
		chunkSize = randomInRange(251904, maxChunkSize, c.random)
	}
	return chunkSize
}

func randomInRange(min, max int, random *xrand.Rand) int {
	return random.Intn(max - min + 1)
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

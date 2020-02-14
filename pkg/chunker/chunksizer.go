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

// --------------------------- CUMULATIVE distribution chunk sizer ---------------------------
type slab struct {
	start int
	end   int
}

type cumulativeChunkSizer struct {
	random *xrand.Rand
	slabs  []slab
}

func newCumulativeChunkSizer() ChunkSizer {
	randomSource := xrand.NewSource(uint64(time.Now().UnixNano()))
	return &cumulativeChunkSizer{
		random: xrand.New(randomSource),
		slabs: []slab{
			slab{start: 96, end: 120},
			slab{start: 120, end: 152},
			slab{start: 152, end: 192},
			slab{start: 192, end: 304},
			slab{start: 304, end: 480},
			slab{start: 480, end: 752},
			slab{start: 752, end: 944},
			slab{start: 944, end: 1228},
			slab{start: 1228, end: 1433},
			slab{start: 1433, end: 1843},
			slab{start: 1843, end: 2355},
			slab{start: 2355, end: 2867},
			slab{start: 2867, end: 3584},
			slab{start: 3584, end: 4505},
			slab{start: 4505, end: 5632},
			slab{start: 5632, end: 7065},
			slab{start: 7065, end: 8908},
			slab{start: 8908, end: 11059},
			slab{start: 11059, end: 13926},
			slab{start: 13926, end: 17305},
			slab{start: 17305, end: 21708},
			slab{start: 21708, end: 27136},
			slab{start: 27136, end: 33894},
			slab{start: 33894, end: 42393},
			slab{start: 42393, end: 52940},
			slab{start: 52940, end: 66252},
			slab{start: 66252, end: 82841},
			slab{start: 82841, end: 103526},
			slab{start: 103526, end: 129331},
			slab{start: 129331, end: 161689},
			slab{start: 161689, end: 202137},
			slab{start: 202137, end: 252723},
			slab{start: 252723, end: 315904},
			slab{start: 315904, end: 394854},
			slab{start: 394854, end: 524288},
			slab{start: 524288, end: 1048576},
		},
	}
}

func (c *cumulativeChunkSizer) NextChunkSize() int {
	prob := c.random.Float64()
	var cumulativeProb float64
	for _, slab := range c.slabs {
		cumulativeProb += (1 / float64(len(c.slabs)))
		if cumulativeProb > prob {
			return randomInRange(slab.start, slab.end, c.random)
		}
	}
	return randomInRange(524288, 1048576, c.random)
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

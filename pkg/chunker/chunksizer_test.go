package chunker

import (
	"testing"
)

func TestChunkSizerReturnsValueWithinRange(t *testing.T) {
	sizer := newRandomChunkSizer()
	for i := 0; i < 100; i++ {
		size := sizer.NextChunkSize()
		if size < minChunkSize || size > maxChunkSize {
			t.Errorf("got size %d but expected size between %d and %d", size, minChunkSize, maxChunkSize)
		}
	}
}

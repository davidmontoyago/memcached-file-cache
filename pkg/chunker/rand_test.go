package chunker

import (
	"testing"
)

func TestChunkSizerReturnsValueWithinRange(t *testing.T) {
	sizer := newChunkSizer()
	for i := 0; i < 100; i++ {
		size := sizer.new()
		if size < minChunkSize || size > maxChunkSize {
			t.Errorf("got size %d but expected size between %d and %d", size, minChunkSize, maxChunkSize)
		}
	}
}

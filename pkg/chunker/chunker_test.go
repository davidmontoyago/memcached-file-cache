package chunker

import (
	"testing"
)

func TestSplitsFileIntoRandomSizedChunks(t *testing.T) {
	chunker := &Chunker{}
	// a 1329 bytes file
	file := []byte(
		"test test test test test test test test test test test test test test test test test test test " +
			"test test test test test test test test test test test test test test test test test test test " +
			"test test test test test test test test test test test test test test test test test test test " +
			"test test test test test test test test test test test test test test test test test test test " +
			"test test test test test test test test test test test test test test test test test test test " +
			"test test test test test test test test test test test test test test test test test test test " +
			"test test test test test test test test test test test test test test test test test test test " +
			"test test test test test test test test test test test test test test test test test test test " +
			"test test test test test test test test test test test test test test test test test test test " +
			"test test test test test test test test test test test test test test test test test test test " +
			"test test test test test test test test test test test test test test test test test test test " +
			"test test test test test test test test test test test test test test test test test test test " +
			"test test test test test test test test test test test test test test test test test test test " +
			"test test test test test test test test test test test test test test test test test test test")
	chunkedFile := chunker.Split(file)

	totalParts := len(chunkedFile.Parts)
	if totalParts < 2 {
		t.Errorf("got %d parts but expected more than 2 parts", totalParts)
	}
	var chunksSizeSum int
	for _, chunk := range chunkedFile.Parts {
		chunksSizeSum += len(chunk.Bytes)
	}
	if chunksSizeSum != 1329 {
		t.Errorf("got %d from adding up chunks but expected %d", chunksSizeSum, 1329)
	}
}

func TestChunkSizerReturnsValueWithinRange(t *testing.T) {
	sizer := newChunkSizer()
	for i := 0; i < 100; i++ {
		size := sizer.new()
		if size < minChunkSize || size > maxChunkSize {
			t.Errorf("got size %d but expected size between %d and %d", size, minChunkSize, maxChunkSize)
		}
	}
}

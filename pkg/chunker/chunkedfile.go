package chunker

import (
	"crypto/md5"
	"fmt"
)

// ChunkedFile represents a file that can be broken into randomly sized chunks or assembled from chunks
type ChunkedFile struct {
	file     []byte
	checksum string
	parts    []*Chunk
}

// NewFromFile initializes a chunked file from a whole file
func NewFromFile(file []byte) *ChunkedFile {
	return &ChunkedFile{file: file}
}

// NewFromChunks creates a chunked file from its parts
func NewFromChunks(parts []*Chunk) *ChunkedFile {
	chunkedFile := NewFromFile(assemble(parts))
	return chunkedFile
}

// Chunks returns a file's chunks - lazy
func (c *ChunkedFile) Chunks() []*Chunk {
	if len(c.parts) == 0 {
		c.split()
	}
	return c.parts
}

// File returns as file
func (c *ChunkedFile) File() []byte {
	return c.file
}

// Validate the file's content checksum against a given checksum
func (c *ChunkedFile) Validate(checksum string) error {
	if checksum != c.Checksum() {
		return fmt.Errorf("content checksum %s does not match %s", c.Checksum(), checksum)
	}
	return nil
}

// Checksum returns a chunked file's checksum
func (c *ChunkedFile) Checksum() string {
	if c.checksum == "" {
		c.checksum = fmt.Sprintf("%x", md5.Sum(c.file))
	}
	return c.checksum
}

// re-assembles a file from its chunks
func assemble(chunks []*Chunk) []byte {
	var file []byte
	for _, chunk := range chunks {
		file = append(file, chunk.Bytes()...)
	}
	return file
}

// Split array of bytes into random sized chunks between 96 and 1024 Kbytes
// chunks between 96 and 1024Kbytes will be randomly distributed across all slabs preventing contention on a single slab
func (c *ChunkedFile) split() {
	var parts []*Chunk

	chunkSizer := newCumulativeChunkSizer()
	var offset int
	for offset < len(c.file) {
		chunkSize := nextChunkSize(len(c.file), offset, chunkSizer)
		parts = append(parts, newChunk(c.file, offset, chunkSize))
		offset += chunkSize
	}

	checksum := md5.Sum(c.file)
	c.checksum = fmt.Sprintf("%x", checksum)
	c.parts = parts
}

func newChunk(file []byte, offset, chunkSize int) *Chunk {
	chunk := file[offset : offset+chunkSize]
	return &Chunk{bytes: chunk}
}

func nextChunkSize(fileSize, offset int, sizer ChunkSizer) int {
	bytesLeft := fileSize - offset
	chunkSize := sizer.NextChunkSize()
	if chunkSize > bytesLeft {
		chunkSize = bytesLeft
	}
	return chunkSize
}

package chunker

import (
	"crypto/md5"
	"fmt"
)

// ChunkedFile is a chunked file with a checksum for consistency verification
type ChunkedFile struct {
	file     []byte
	checksum string
	parts    []*Chunk
}

// NewFromFile initializes a chunked file from a whole file
func NewFromFile(file []byte) *ChunkedFile {
	return &ChunkedFile{file: file}
}

// NewFromChunks creates a chunked files from its parts
func NewFromChunks(checksum string, parts []*Chunk) *ChunkedFile {
	return &ChunkedFile{checksum: checksum, parts: parts}
}

// Checksum returns a chunked file's checksum
func (c *ChunkedFile) Checksum() string {
	return c.checksum
}

// Chunks returns a file chunks
func (c *ChunkedFile) Chunks() []*Chunk {
	if len(c.parts) == 0 {
		c.split()
	}
	return c.parts
}

// Assemble re-assembles a file from its chunks
func (c *ChunkedFile) Assemble() []byte {
	c.file = nil
	for _, chunk := range c.Chunks() {
		c.file = append(c.file, chunk.Bytes()...)
	}
	return c.file
}

// Validate ensures a chunked file checksum matches its assembled parts
func (c *ChunkedFile) Validate() error {
	fileChecksum := fmt.Sprintf("%x", md5.Sum(c.file))
	if fileChecksum != c.Checksum() {
		return fmt.Errorf("chunked file checksum %s does not match its content: %s", c.Checksum(), fileChecksum)
	}
	return nil
}

// Split array of bytes into random sized chunks between 96 and 1024 Kbytes
// chunks between 96 and 1024Kbytes will be randomly distributed across all slabs preventing contention on a single slab
func (c *ChunkedFile) split() {
	var parts []*Chunk

	chunkSizer := newChunkSizer()
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

func nextChunkSize(fileSize, offset int, sizer *chunkSizer) int {
	bytesLeft := fileSize - offset
	chunkSize := sizer.new()
	if chunkSize > bytesLeft {
		chunkSize = bytesLeft
	}
	return chunkSize
}

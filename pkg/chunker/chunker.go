package chunker

import (
	"crypto/md5"
	"fmt"
	"log"
)

// Chunker allows splitting files into random sized chunks
type Chunker struct {
}

// Chunk if a part of a file
type Chunk struct {
	bytes []byte
}

// Bytes returns a chunk's bytes
func (c *Chunk) Bytes() []byte {
	return c.bytes
}

// Chunked is a chunked file with a checksum for consistency verification
type Chunked struct {
	checksum string
	parts    []*Chunk
}

// Checksum returns a chunked file's checksum
func (c *Chunked) Checksum() string {
	return c.checksum
}

// Chunks returns a file chunks
func (c *Chunked) Chunks() []*Chunk {
	return c.parts
}

// Split array of bytes into random sized chunks between 96 and 1024 Kbytes
// chunks between 96 and 1024Kbytes will be randomly distributed across all slabs preventing contention on a single slab
func (c *Chunker) Split(file []byte) *Chunked {
	log.Println("file size is", len(file), "bytes")
	var parts []*Chunk

	chunkSizer := newChunkSizer()
	var offset int
	for offset < len(file) {
		chunkSize := nextChunkSize(len(file), offset, chunkSizer)
		parts = append(parts, newChunk(file, offset, chunkSize))
		offset += chunkSize
	}

	checksum := md5.Sum(file)
	return &Chunked{
		checksum: fmt.Sprintf("%x", checksum),
		parts:    parts,
	}
}

// Assemble re-assembles a file from its chunks
func (c *Chunker) Assemble(chunkedFile *Chunked) []byte {
	var file []byte
	for _, chunk := range chunkedFile.Chunks() {
		file = append(file, chunk.Bytes()...)
	}
	return file
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

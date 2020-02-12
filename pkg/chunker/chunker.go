package chunker

import (
	"crypto/sha256"
	"log"
)

// Chunker allows splitting files into random sized chunks
type Chunker struct {
}

// Chunk if a part of a file
type Chunk struct {
	Bytes []byte
}

// Chunked is a chunked file with a checksum for consistency verification
type Chunked struct {
	Checksum [sha256.Size]byte
	Parts    []*Chunk
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
	return &Chunked{Parts: parts}
}

func newChunk(file []byte, offset, chunkSize int) *Chunk {
	chunk := file[offset : offset+chunkSize]
	return &Chunk{
		Bytes: chunk,
	}
}

func nextChunkSize(fileSize, offset int, sizer *chunkSizer) int {
	bytesLeft := fileSize - offset
	chunkSize := sizer.new()
	if chunkSize > bytesLeft {
		chunkSize = bytesLeft
	}
	log.Println("next chunk size is", chunkSize)
	return chunkSize
}

// Assemble re-assembles a file from its chunks
func (c *Chunker) Assemble(chunkedFile *Chunked) []byte {
	var file []byte
	for _, chunk := range chunkedFile.Parts {
		file = append(file, chunk.Bytes...)
	}
	return file
}

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

// Split array of bytes into random sized chunks between 96 and 1024 bytes
// a chunk of at least 96 bytes will ensure it fits onto any page of the first slab
// chunks between and 96 and 1024 will be distributed across all slabs reducing contention on a single slab
func (c *Chunker) Split(file []byte) *Chunked {
	log.Println("file size is", len(file), "bytes")
	var parts []*Chunk

	chunkSizer := newChunkSizer()
	var i int
	for i < len(file) {
		bytesLeft := len(file) - i
		chunkSize := chunkSizer.new()
		if chunkSize > bytesLeft {
			chunkSize = bytesLeft
		}
		log.Println("next chunk size is", chunkSize)

		chunk := file[i : i+chunkSize]
		parts = append(parts, &Chunk{
			Bytes: chunk,
		})
		i += chunkSize
	}

	return &Chunked{
		Parts: parts,
	}
}

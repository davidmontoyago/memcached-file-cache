package chunker

// Chunk if a part of a file
type Chunk struct {
	bytes []byte
}

// NewChunk inits a chunk from an array of bytes
func NewChunk(bytes []byte) *Chunk {
	return &Chunk{bytes: bytes}
}

// Bytes returns a chunk's bytes
func (c *Chunk) Bytes() []byte {
	return c.bytes
}

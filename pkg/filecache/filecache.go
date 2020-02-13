package filecache

import (
	"fmt"
	"strings"

	"github.com/bradfitz/gomemcache/memcache"
	"github.com/davidmontoyago/interview-davidmontoyago-d660952eff664d8bac96c9124d7f8582/pkg/chunker"
)

// FileCache interfaces with cache client and chunker apis for spliting and storing/fetching files
type FileCache struct {
	memcache memcachedClient
}

// NewFileCache initializes a FileCache with a memcache client
func NewFileCache(memcache memcachedClient) *FileCache {
	return &FileCache{
		memcache: memcache,
	}
}

// Put splits and places a file in the cache
func (f *FileCache) Put(file []byte) error {
	chunkedFile := chunker.NewFromFile(file)
	totalChunks := len(chunkedFile.Chunks())

	var sb strings.Builder
	fileChunksByKey := make(map[string]*chunker.Chunk, totalChunks)
	for count, chunk := range chunkedFile.Chunks() {
		key := fmt.Sprintf("%s-part-%d-of-%d", chunkedFile.Checksum(), count, totalChunks)
		sb.WriteString(key)
		sb.WriteString(",")
		fileChunksByKey[key] = chunk
	}
	keys := strings.TrimSuffix(sb.String(), ",")

	// store file unique id and comma separated list of all its chunks' ids
	fileKeys := &memcache.Item{Key: chunkedFile.Checksum(), Value: []byte(keys)}
	if err := f.memcache.Set(fileKeys); err != nil {
		return err
	}

	// store file chunks
	for key, chunk := range fileChunksByKey {
		chunkItem := &memcache.Item{Key: key, Value: chunk.Bytes()}
		if err := f.memcache.Set(chunkItem); err != nil {
			return err
		}
	}
	return nil
}

// Get fetches all file parts and returns the assembled file
func (f *FileCache) Get(key string) ([]byte, error) {
	fileKeys, err := f.memcache.Get(key)
	if err != nil {
		return nil, err
	}

	var parts []*chunker.Chunk
	chunksKeys := string(fileKeys.Value)
	for _, chunkKey := range strings.Split(chunksKeys, ",") {
		fileChunk, err := f.memcache.Get(chunkKey)
		if err != nil {
			return nil, err
		}
		parts = append(parts, chunker.NewChunk(fileChunk.Value))
	}

	chunkedFile := chunker.NewFromChunks(key, parts)
	return chunkedFile.Assemble(), nil
}

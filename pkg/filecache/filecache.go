package filecache

import (
	"fmt"
	"strings"

	"github.com/bradfitz/gomemcache/memcache"
	"github.com/davidmontoyago/interview-davidmontoyago-d660952eff664d8bac96c9124d7f8582/pkg/chunker"
)

// FileCache interfaces with cache client and chunker apis for spliting and storing/fetching files
type FileCache struct {
	chunker  *chunker.Chunker
	memcache memcachedClient
}

// NewFileCache initializes a FileCache with a memcache client
func NewFileCache(memcache memcachedClient) *FileCache {
	return &FileCache{
		chunker:  &chunker.Chunker{},
		memcache: memcache,
	}
}

// Put splits and places a file in the cache
func (f *FileCache) Put(file []byte) {
	chunkedFile := f.chunker.Split(file)
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
	f.memcache.Set(fileKeys)

	// store file chunks
	for key, chunk := range fileChunksByKey {
		chunkItem := &memcache.Item{Key: key, Value: chunk.Bytes()}
		f.memcache.Set(chunkItem)
	}
}

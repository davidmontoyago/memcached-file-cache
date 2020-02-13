package filecache

import (
	"fmt"
	"log"
	"strings"

	"github.com/bradfitz/gomemcache/memcache"
	"github.com/davidmontoyago/interview-davidmontoyago-d660952eff664d8bac96c9124d7f8582/pkg/chunker"
	"github.com/pkg/errors"
)

// MaxFileSize 50MB
const MaxFileSize = 52428800

// FileCache interfaces with cache client and chunker apis for spliting and storing/fetching files
type FileCache struct {
	memcache memcachedClient
}

// New initializes a FileCache with a memcache client
func New(memcache memcachedClient) *FileCache {
	return &FileCache{
		memcache: memcache,
	}
}

// Put splits and places a file in the cache
func (f *FileCache) Put(file []byte) (string, error) {
	if err := checkFileSize(file); err != nil {
		return "", err
	}

	chunkedFile := chunker.NewFromFile(file)

	if f.exists(chunkedFile.Checksum()) {
		return chunkedFile.Checksum(), nil
	}

	chunksByKey, keys := getChunksByKey(chunkedFile)

	if err := f.putFileKeys(chunkedFile.Checksum(), keys); err != nil {
		return "", err
	}
	if err := f.putFileChunks(chunksByKey); err != nil {
		log.Println(f.memcache.Delete(chunkedFile.Checksum()))
		return "", err
	}
	return chunkedFile.Checksum(), nil
}

// Get fetches all file parts and returns the assembled file
func (f *FileCache) Get(checksum string) ([]byte, error) {
	var chunksKeys []string
	var err error
	if chunksKeys, err = f.getFileChunksKeys(checksum); err != nil {
		return nil, err
	}

	var chunks []*chunker.Chunk
	if chunks, err = f.getFileChunks(chunksKeys); err != nil {
		return nil, err
	}

	chunkedFile := chunker.NewFromChunks(chunks)

	if err := chunkedFile.Validate(checksum); err != nil {
		return nil, errors.Wrap(err, "failed to validate file")
	}

	return chunkedFile.File(), nil
}

// store file unique id and comma separated list of all its chunks' ids
func (f *FileCache) putFileKeys(fileKey string, chunksKeys string) error {
	fileKeys := &memcache.Item{Key: fileKey, Value: []byte(chunksKeys)}
	if err := f.memcache.Set(fileKeys); err != nil {
		return errors.Wrap(err, "failed to put file keys")
	}
	return nil
}

// store file chunks
func (f *FileCache) putFileChunks(chunksByKey map[string]*chunker.Chunk) error {
	for key, chunk := range chunksByKey {
		chunkItem := &memcache.Item{Key: key, Value: chunk.Bytes()}
		if err := f.memcache.Set(chunkItem); err != nil {
			return errors.Wrap(err, "failed to put file's chunks")
		}
	}
	return nil
}

// get the keys for all the file's chunks
func (f *FileCache) getFileChunksKeys(checksum string) ([]string, error) {
	var keys []string
	fileKeys, err := f.memcache.Get(checksum)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get file keys")
	}
	chunksKeys := string(fileKeys.Value)
	if chunksKeys != "" {
		keys = strings.Split(chunksKeys, ",")
	}
	return keys, nil
}

// get all chunks given a set of keys
func (f *FileCache) getFileChunks(chunksKeys []string) ([]*chunker.Chunk, error) {
	var parts []*chunker.Chunk
	for _, chunkKey := range chunksKeys {
		fileChunk, err := f.memcache.Get(chunkKey)
		if err != nil {
			return nil, errors.Wrap(err, "failed to get file chunk")
		}
		parts = append(parts, chunker.NewChunk(fileChunk.Value))
	}
	return parts, nil
}

// map every file chunk to its respective key; also returns comma separated list of keys ready for storing in memcached
func getChunksByKey(chunkedFile *chunker.ChunkedFile) (map[string]*chunker.Chunk, string) {
	var sb strings.Builder

	totalChunks := len(chunkedFile.Chunks())
	chunksByKey := make(map[string]*chunker.Chunk, totalChunks)
	for count, chunk := range chunkedFile.Chunks() {
		key := fmt.Sprintf("%s-part-%d-of-%d", chunkedFile.Checksum(), count, totalChunks)
		chunksByKey[key] = chunk

		sb.WriteString(key)
		sb.WriteString(",")
	}
	return chunksByKey, strings.TrimSuffix(sb.String(), ",")
}

func checkFileSize(file []byte) error {
	fileSize := len(file)
	if fileSize > MaxFileSize {
		return fmt.Errorf("file size %d exceeds the max allowed %d by %d bytes", fileSize, MaxFileSize, fileSize-MaxFileSize)
	}
	return nil
}

func (f *FileCache) exists(checksum string) bool {
	item, err := f.memcache.Get(checksum)
	if err != nil {
		return false
	}
	return item != nil
}

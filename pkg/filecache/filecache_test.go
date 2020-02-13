package filecache

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/bradfitz/gomemcache/memcache"
)

type mockMemcachedClient struct {
	storedItems map[string][]byte
}

func (c *mockMemcachedClient) Set(item *memcache.Item) error {
	c.storedItems[item.Key] = item.Value
	return nil
}

func TestPutStoresFileKeyWithChunksKeysAsValue(t *testing.T) {
	memcached := &mockMemcachedClient{storedItems: make(map[string][]byte)}
	fileCache := NewFileCache(memcached)

	f, err := os.Open("../chunker/fixture/file.dat")
	if err != nil {
		t.Error(err)
	}
	file, err := ioutil.ReadAll(f)
	if err != nil {
		t.Error(err)
	}
	fileCache.Put(file)

	fileKey := "91388263e7c545ebea3952fb2637dffa"
	var val []byte
	var ok bool
	if val, ok = memcached.storedItems[fileKey]; !ok {
		t.Errorf("expected file with key %s to be stored but found none", fileKey)
	}
	chunksKeys := string(val)
	for count, chunkKey := range strings.Split(chunksKeys, ",") {
		if val, ok = memcached.storedItems[chunkKey]; !ok {
			t.Errorf("expected file chunk key %s to be stored but found none", chunkKey)
		}
		expectedKey := fmt.Sprintf("%s-part-%d-of-%d", fileKey, count, len(memcached.storedItems)-1)
		if chunkKey != expectedKey {
			t.Errorf("expected file chunk key to be %s but got %s", expectedKey, chunkKey)
		}
	}
}

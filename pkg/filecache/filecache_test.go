package filecache

import (
	"bytes"
	"crypto/rand"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/pkg/errors"

	"github.com/bradfitz/gomemcache/memcache"
)

type mockMemcachedClient struct {
	storedItems    map[string][]byte
	getTimesCalled int
}

func (c *mockMemcachedClient) Set(item *memcache.Item) error {
	c.storedItems[item.Key] = item.Value
	return nil
}

func (c *mockMemcachedClient) Get(key string) (item *memcache.Item, err error) {
	c.getTimesCalled++
	val, ok := c.storedItems[key]
	if !ok {
		return nil, errors.New("cache miss")
	}
	return &memcache.Item{Value: val}, nil
}

func (c *mockMemcachedClient) Delete(key string) error {
	return nil
}

func TestPutStoresFileKeyWithChunksKeysAsValue(t *testing.T) {
	expectedFileKey := "91388263e7c545ebea3952fb2637dffa"
	memcached := &mockMemcachedClient{storedItems: make(map[string][]byte)}
	fileCache := New(memcached)

	f, err := os.Open("../chunker/fixture/file.dat")
	if err != nil {
		t.Error(err)
	}
	file, err := ioutil.ReadAll(f)
	if err != nil {
		t.Error(err)
	}
	key, err := fileCache.Put(file)
	if err != nil {
		t.Error(err)
	}
	if key != expectedFileKey {
		t.Errorf("got file key %s but expected %s", key, expectedFileKey)
	}

	var val []byte
	var ok bool
	if val, ok = memcached.storedItems[expectedFileKey]; !ok {
		t.Errorf("expected file with key %s to be stored but found none", expectedFileKey)
	}
	chunksKeys := string(val)
	for count, chunkKey := range strings.Split(chunksKeys, ",") {
		if val, ok = memcached.storedItems[chunkKey]; !ok {
			t.Errorf("expected file chunk key %s to be stored but found none", chunkKey)
		}
		expectedKey := fmt.Sprintf("%s-part-%d-of-%d", expectedFileKey, count, len(memcached.storedItems)-1)
		if chunkKey != expectedKey {
			t.Errorf("expected file chunk key to be %s but got %s", expectedKey, chunkKey)
		}
	}
}

func TestGetFetchesChunksAndAssemblesFile(t *testing.T) {
	memcached := &mockMemcachedClient{storedItems: make(map[string][]byte)}
	fileCache := New(memcached)

	f, err := os.Open("../chunker/fixture/file.dat")
	if err != nil {
		t.Error(err)
	}
	file, err := ioutil.ReadAll(f)
	if err != nil {
		t.Error(err)
	}
	_, err = fileCache.Put(file)
	if err != nil {
		t.Error(err)
	}

	assembledFile, err := fileCache.Get("91388263e7c545ebea3952fb2637dffa")
	if err != nil {
		t.Error(err)
	}
	if !bytes.Equal(file, assembledFile) {
		t.Errorf("expected assembled file to contain same bytes as fixture but they differ")
	}
}

func TestPutRejectsFilesGreaterThan50MB(t *testing.T) {
	bigFile := make([]byte, MaxFileSize+1)
	rand.Read(bigFile)

	memcached := &mockMemcachedClient{storedItems: make(map[string][]byte)}
	fileCache := New(memcached)

	_, err := fileCache.Put(bigFile)
	if err == nil {
		t.Error("expected error but got nil")
	}
}

func TestGetCanHandleAnEmptyFile(t *testing.T) {
	memcached := &mockMemcachedClient{storedItems: make(map[string][]byte)}
	fileCache := New(memcached)

	f, err := os.Open("../chunker/fixture/empty.dat")
	if err != nil {
		t.Error(err)
	}
	file, err := ioutil.ReadAll(f)
	if err != nil {
		t.Error(err)
	}
	_, err = fileCache.Put(file)
	if err != nil {
		t.Error(err)
	}

	_, err = fileCache.Get("d41d8cd98f00b204e9800998ecf8427e")
	if err != nil {
		t.Error(err)
	}

	if memcached.getTimesCalled > 2 {
		t.Errorf("expected memcached Get to be called twice but was called %d times", memcached.getTimesCalled)
	}
}

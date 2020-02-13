package filecache

import "github.com/bradfitz/gomemcache/memcache"

type memcachedClient interface {
	Set(item *memcache.Item) error
}

package server

import (
	"time"

	"github.com/dgraph-io/ristretto"
	"github.com/prometheus/common/log"
	"github.com/vmihailenco/msgpack/v5"
)

var (
	cache        *ristretto.Cache
	itemLifetime time.Duration
	cacheOff     bool
)

// Cache for the system
func initCache(settings ConfigSchema) (err error) {
	cacheOff = !settings.CacheEnabled
	if cacheOff {
		log.Infoln("cache is disabled")
		return
	}
	log.Infoln("cache is enabled")
	cache, err = ristretto.NewCache(&ristretto.Config{
		NumCounters: int64(settings.CacheMaxSize) * 10,   // cost x10
		MaxCost:     int64(settings.CacheMaxSize) * 1024, // in megabytes
		BufferItems: 64,                                  // number of keys per Get buffer.
	})
	itemLifetime = time.Duration(settings.CacheLifetime) * time.Second
	return
}

// set store an element in the cache
func set(k, v interface{}) (err error) {
	if cacheOff {
		return
	}
	// serialize the data
	data, err := msgpack.Marshal(v)
	if err != nil {
		return
	}
	// the cost is the size of the item
	cost := int64(len(data))
	// set the data with ttl
	cache.SetWithTTL(k, data, cost, itemLifetime)
	return
}

// Get from the cache
func get(k, v interface{}) (found bool) {
	if cacheOff {
		return
	}
	// value
	data, found := cache.Get(k)
	if !found {
		return
	}
	// serialize the data
	if err := msgpack.Unmarshal(data.([]byte), v); err != nil {
		found = false
	}
	return
}

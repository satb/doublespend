package cache

/*

The cache is a map of string to cache items
For the purposes of doublespend, we treat the key of the map as the address
to monitor and the cacheItems as a list of structs that are the transactions
that were performed by the key (address).

*/

import (
	"errors"
	"log"
	"math/big"
	"sync"
	"time"
)

type Cache struct {
	RetentionPeriod uint64
	store           map[string][]CacheItem
	rwm             sync.RWMutex
}

type CacheItem struct {
	Id          string
	From        string
	To          string
	Amount      big.Int
	BlockNum    int64
	Time        uint64
	DoubleSpend bool
}

type updateFn func(CacheItem) CacheItem

const CACHE_ITEM_SIZE = 0
const DEFAULT_RETAIN_PERIOD uint64 = 24 * 60 * 60 * 1000

func New() *Cache {
	return &Cache{
		RetentionPeriod: DEFAULT_RETAIN_PERIOD,
		store:           make(map[string][]CacheItem),
	}
}

/*
Returns a COPY of the items associated with the cache key.
Since the items are a slice, the slice is deep copied and returned.
The reason for doing so is because the item can be purged while the consumer of this
item might be working on it. Working with copies makes it safe to return the item
to the clients that are free to do anything they want thereafter with the item
*/
func (cache *Cache) Get(key string) []CacheItem {
	cache.rwm.RLock()
	defer cache.rwm.RUnlock()
	return append([]CacheItem(nil), cache.store[key]...)
}

func (cache *Cache) AddItem(key string, item CacheItem) {
	cache.rwm.Lock()
	defer cache.rwm.Unlock()
	list := cache.store[key]
	if list == nil {
		list = make([]CacheItem, CACHE_ITEM_SIZE)
	}
	cache.store[key] = append(list, item)
}

func (cache *Cache) Del(key string) {
	cache.rwm.Lock()
	defer cache.rwm.Unlock()
	if _, ok := cache.store[key]; ok {
		delete(cache.store, key)
	}
}

func (cache *Cache) UpdateItem(key string, id string, fn updateFn) {
	cache.rwm.Lock()
	defer cache.rwm.Unlock()
	items := cache.store[key]
	var it *CacheItem
	for _, item := range items {
		if item.Id == id {
			it = &item
			break
		}
	}
	if it != nil {
		updatedItem := fn(*it)
		newList := make([]CacheItem, 0)
		for _, item := range items {
			if item.Id == id {
				newList = append(newList, updatedItem)
			} else {
				newList = append(newList, item)
			}
		}
		cache.store[key] = newList
	} else {
		log.Println("Could not find transaction with key=", key, " and id=", id)
	}
}

/*
Returns the item associated with the given key and the id.
The key is the index into the map and the id is the identifier for the item
being requested.
*/
func (cache *Cache) GetItem(key string, id string) (item CacheItem, err error) {
	cache.rwm.RLock()
	defer cache.rwm.RUnlock()
	items := cache.store[key]
	for _, item := range items {
		if item.Id == id {
			return item, nil
		}
	}
	return CacheItem{}, errors.New("Not found")
}

/**
Builds a list of all transactions that can be purged and purges them from the cache
We have a time based purge that allows us to purge the transactions that are over
the given retain period.
*/

func (cache *Cache) Purge(key string) {
	now := uint64(time.Now().Unix())
	items := cache.Get(key)
	if items != nil {
		list := make([]string, 0)
		for _, item := range items {
			if (now - item.Time) > cache.RetentionPeriod {
				list = append(list, item.Id)
			}
		}
		if len(items) > 0 {
			cache.deleteItems(key, list)
		}
	}
}

/**
Allocate a new cache items buffer and get rid of the old one for gc
*/
func (cache *Cache) deleteItems(key string, ids []string) {
	cache.rwm.Lock()
	defer cache.rwm.Unlock()
	if items, ok := cache.store[key]; ok {
		newItems := make([]CacheItem, CACHE_ITEM_SIZE)
		for _, item := range items {
			found := false
			for _, id := range ids {
				if id == item.Id {
					found = true
					break
				}
			}
			if !found {
				newItems = append(newItems, item)
			}
		}
		cache.store[key] = newItems
	}
}

/**
Clears the cache completely. Useful for tests but not so much for real use
*/
func (cache *Cache) clear() {
	cache.rwm.Lock()
	defer cache.rwm.Unlock()
	cache.store = make(map[string][]CacheItem)
}
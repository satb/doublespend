package cache

import (
	"math/big"
	"testing"
	"time"
)

var c = New()

func setup() {
	c.RetentionPeriod = DEFAULT_RETAIN_PERIOD
	c.clear()
}

func TestAddItem(t *testing.T) {
	setup()
	key := "0x00"
	item := CacheItem{
		Id:       "id1",
		From:     "from1",
		To:       "to1",
		Amount:   *new(big.Int).SetInt64(10),
		BlockNum: 10,
	}
	c.AddItem(key, item)
	if c.Get(key) == nil {
		t.Log("Putting item in cache failed")
		t.Fail()
	}
}

func TestDel(t *testing.T) {
	setup()
	key := "0x00"
	item := CacheItem{
		Id:       "id1",
		From:     "from1",
		To:       "to1",
		Amount:   *new(big.Int).SetInt64(10),
		BlockNum: 10,
	}
	c.AddItem(key, item)
	c.Del(key)
	if c.Get(key) != nil {
		t.Log("Deleting item in cache failed")
		t.Fail()
	}
}

func TestAddItems(t *testing.T) {
	setup()
	key := "0x00"
	item1 := CacheItem{
		Id:       "id1",
		From:     "from1",
		To:       "to1",
		Amount:   *new(big.Int).SetInt64(10),
		BlockNum: 10,
	}
	item2 := CacheItem{
		Id:       "id2",
		From:     "from2",
		To:       "to2",
		Amount:   *new(big.Int).SetInt64(12),
		BlockNum: 11,
	}
	c.AddItem(key, item1)
	c.AddItem(key, item2)
	items := c.Get(key)
	if items == nil || len(items) != CACHE_ITEM_SIZE+2 {
		t.Log("AddItems failed")
		t.Fail()
	}
}

func TestDeleteItems(t *testing.T) {
	setup()
	key := "0x00"
	item1 := CacheItem{
		Id:       "id1",
		From:     "from1",
		To:       "to1",
		Amount:   *new(big.Int).SetInt64(10),
		BlockNum: 10,
	}
	item2 := CacheItem{
		Id:       "id2",
		From:     "from2",
		To:       "to2",
		Amount:   *new(big.Int).SetInt64(12),
		BlockNum: 11,
	}
	c.AddItem(key, item1)
	c.AddItem(key, item2)
	c.deleteItems(key, []string{item1.Id})
	items := c.Get(key)
	if items == nil || len(items) != CACHE_ITEM_SIZE+1 || items[0].Id != "id2" {
		t.Log("Failed to delete items")
		t.Fail()
	}
}

func TestUpdateItem(t *testing.T) {
	setup()
	key := "0x00"
	item1 := CacheItem{
		Id:       "id1",
		From:     "from1",
		To:       "to1",
		Amount:   *new(big.Int).SetInt64(10),
		BlockNum: 10,
	}
	item2 := CacheItem{
		Id:       "id2",
		From:     "from2",
		To:       "to2",
		Amount:   *new(big.Int).SetInt64(12),
		BlockNum: 11,
	}
	c.AddItem(key, item1)
	c.AddItem(key, item2)
	c.UpdateItem(key, item1.Id, func(ci CacheItem) CacheItem {
		ci.Id = "id3"
		ci.From = "from3"
		return ci
	})
	c.UpdateItem(key, item2.Id, func(ci CacheItem) CacheItem {
		ci.Id = "id4"
		ci.From = "from4"
		return ci
	})
	it1 := c.Get(key)[0]
	it2 := c.Get(key)[1]
	if it1.Id != "id3" || it1.From != "from3" || it2.Id != "id4" || it2.From != "from4" {
		t.Log("Failed to update item")
		t.Fail()
	}
}

func TestPurge(t *testing.T) {
	setup()
	key := "0x00"
	item := CacheItem{
		Id:       "id1",
		From:     "from1",
		To:       "to1",
		Amount:   *new(big.Int).SetInt64(10),
		BlockNum: 10,
	}
	c.AddItem(key, item)
	c.RetentionPeriod = 10
	time.Sleep(50)
	c.Purge(key)
	if c.Get(key) != nil {
		t.Log("Failed to purge item")
		t.Fail()
	}
}

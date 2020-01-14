package cache

import (
	"math/big"
	"testing"
	"time"
)

var c = New()

func setup() {
	c.RetentionPeriod = DefaultRetainPeriod
	c.clear()
}

func TestGetWithEmptyCache(t *testing.T) {
	setup()
	key := "0x00"
	if c.Get(key) != nil {
		t.Log("Getting item from empty cache failed")
		t.Fail()
	}
}

func TestAddItem(t *testing.T) {
	setup()
	key := "0x00"
	item := Item{
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
	item := Item{
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
	item1 := Item{
		Id:       "id1",
		From:     "from1",
		To:       "to1",
		Amount:   *new(big.Int).SetInt64(10),
		BlockNum: 10,
	}
	item2 := Item{
		Id:       "id2",
		From:     "from2",
		To:       "to2",
		Amount:   *new(big.Int).SetInt64(12),
		BlockNum: 11,
	}
	c.AddItem(key, item1)
	c.AddItem(key, item2)
	items := c.Get(key)
	if items == nil || len(items) != ItemSize+2 {
		t.Log("AddItems failed")
		t.Fail()
	}
}

func TestDeleteItems(t *testing.T) {
	setup()
	key := "0x00"
	item1 := Item{
		Id:       "id1",
		From:     "from1",
		To:       "to1",
		Amount:   *new(big.Int).SetInt64(10),
		BlockNum: 10,
	}
	item2 := Item{
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
	if items == nil || len(items) != ItemSize+1 || items[0].Id != "id2" {
		t.Log("Failed to delete items")
		t.Fail()
	}
}

func TestUpdateItem(t *testing.T) {
	setup()
	key := "0x00"
	item1 := Item{
		Id:       "id1",
		From:     "from1",
		To:       "to1",
		Amount:   *new(big.Int).SetInt64(10),
		BlockNum: 10,
	}
	item2 := Item{
		Id:       "id2",
		From:     "from2",
		To:       "to2",
		Amount:   *new(big.Int).SetInt64(12),
		BlockNum: 11,
	}
	c.AddItem(key, item1)
	c.AddItem(key, item2)

	item3 := item1
	item3.Id = "id3"
	item3.From = "from3"
	c.UpdateItem(key, item3)

	item4 := item1
	item4.Id = "id4"
	item4.From = "from4"
	c.UpdateItem(key, item4)
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
	item := Item{
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

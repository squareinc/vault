package keysutil

import lru "github.com/hashicorp/golang-lru"

type TransitLRU struct {
	lru *lru.TwoQueueCache
}

func NewTransitLRU(size int) *TransitLRU {
	lru, _ := lru.New2Q(size)
	return &TransitLRU{lru: lru}
}

func (c *TransitLRU) Delete(key interface{}) {
	c.lru.Remove(key)
}

func (c *TransitLRU) Load(key interface{}) (value interface{}, ok bool) {
	value, ok = c.lru.Get(key)
	return
}

func (c *TransitLRU) Store(key, value interface{}) {
	c.lru.Add(key, value)
}

func (c *TransitLRU) Purge() {
	c.lru.Purge()
}

func (c *TransitLRU) Len() int {
	return c.lru.Len()
}

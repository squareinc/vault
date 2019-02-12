package keysutil

import (
	"sync"
)

type TransitSyncMap struct {
	syncmap sync.Map
}

func NewTransitSyncMap() *TransitSyncMap {
	return &TransitSyncMap{syncmap: sync.Map{}}
}

func (c *TransitSyncMap) Delete(key interface{}) {
	c.syncmap.Delete(key)
}

func (c *TransitSyncMap) Load(key interface{}) (value interface{}, ok bool) {
	value, ok = c.syncmap.Load(key)
	return
}

func (c *TransitSyncMap) Store(key, value interface{}) {
	c.syncmap.Store(key, value)
}

func (c *TransitSyncMap) Purge() {
	deleteFunc := func(key, value interface{}) bool {
		c.syncmap.Delete(key)
		return true
	}
	c.syncmap.Range(deleteFunc)
}

func (c *TransitSyncMap) Len() int {
	length := 0
	countFunc := func(_, _ interface{}) bool {
		length++
		return true
	}
	c.syncmap.Range(countFunc)
	return length
}

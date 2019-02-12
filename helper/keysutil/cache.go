package keysutil

type CacheType int

const (
	NOTIMPLEMENTED CacheType = iota
	SYNCMAP
	LRU
)

type Cache interface {
	Delete(key interface{})
	Load(key interface{}) (value interface{}, ok bool)
	Store(key, value interface{})
	Purge()
	Len() int
}

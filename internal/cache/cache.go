package cache

import (
	"time"

	gocache "github.com/patrickmn/go-cache"
)

type EntryExists struct{}
func (ee *EntryExists) Error() string {
	return "entry already exists"
}

var (
	entries = gocache.New(5 * time.Minute, 10 * time.Second)
)

func Register(id, address string) error {
	_, ok := entries.Get(id)
	if ok {
		return &EntryExists{}
	}
	return entries.Add(id, address, gocache.DefaultExpiration)
}

func Deregister(id string) {
	entries.Delete(id)
}

func Get(id string) string {
	if entry, ok := entries.Get(id); ok {
		return entry.(string)
	}
	return ""
}

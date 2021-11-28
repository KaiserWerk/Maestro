package cache

import (
	"fmt"
	"sync"
	"time"

	"github.com/KaiserWerk/Maestro/internal/entity"

	"github.com/patrickmn/go-cache"
)

type EntryExists struct{}
func (ee *EntryExists) Error() string {
	return "entry already exists"
}

var (
	mut sync.Mutex
	entries = cache.New(5 * time.Minute, time.Minute)
)

func Register(id, address string) error {
	_, ok := entries.Get(id)
	if ok {
		return &EntryExists{}
	}

	entries.Set(id, &entity.Registrant{
		Id:       id,
		Address:  address,
		LastPing: time.Now(),
	}, cache.DefaultExpiration)
	return nil
}

func Deregister(id string) {
	entries.Delete(id)
}

func Get(id string) (*entity.Registrant, bool) {
	if e, ok := entries.Get(id); ok {
		return e.(*entity.Registrant), true
	}
	return nil, false
}

func Update(id string) error {
	e, ok := entries.Get(id)
	if !ok {
		return fmt.Errorf("entry with id '%s' does not exist", id)
	}

	reg := e.(*entity.Registrant)
	reg.LastPing = time.Now()

	entries.Set(id, reg, cache.DefaultExpiration)
	return nil
}

/*
func Register(id, address string) error {
	mut.Lock()
	defer mut.Unlock()
	_, ok := entries[id]
	if ok {
		return &EntryExists{}
	}
	entries[id] = entity.Registrant{
		Id:       id,
		Address:  address,
		LastPing: time.Now(),
	}
	return nil
}

func Deregister(id string) {
	mut.Lock()
	defer mut.Unlock()
	delete(entries, id)
}

func Get(id string) (entity.Registrant, bool) {
	mut.Lock()
	defer mut.Unlock()
	reg, ok := entries[id]
	return reg, ok
}

func GetAll() []entity.Registrant {
	mut.Lock()
	defer mut.Unlock()
	s := make([]entity.Registrant, 0)
	for _, v := range entries {
		s = append(s, v)
	}
	return s
}
*/
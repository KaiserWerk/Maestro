package cache

import (
	"fmt"
	"github.com/KaiserWerk/Maestro/internal/configuration"
	"time"

	"github.com/KaiserWerk/Maestro/internal/entity"

	"github.com/patrickmn/go-cache"
)

// EntryExists error
type EntryExists struct{}

func (ee *EntryExists) Error() string {
	return "entry already exists"
}

type Cacher interface {
	Register(string, string) error
	Deregister(string) bool
	Get(string) (*entity.Registrant, bool)
	Update(string) error
}

type MaestroCache struct {
	entries *cache.Cache
}

var _ Cacher = &MaestroCache{}

func New(cfg *configuration.AppConfig) *MaestroCache {
	return &MaestroCache{entries: cache.New(cfg.DieAfter, time.Minute)}
}

func (mc *MaestroCache) Register(id, address string) error {
	_, ok := mc.entries.Get(id)
	if ok {
		return &EntryExists{}
	}

	mc.entries.Set(id, &entity.Registrant{
		Id:       id,
		Address:  address,
		LastPing: time.Now(),
	}, cache.DefaultExpiration)
	return nil
}

func (mc *MaestroCache) Deregister(id string) bool {
	if _, ok := mc.entries.Get(id); ok {
		mc.entries.Delete(id)
		return true
	}
	return false
}

func (mc *MaestroCache) Get(id string) (*entity.Registrant, bool) {
	if e, ok := mc.entries.Get(id); ok {
		return e.(*entity.Registrant), true
	}
	return nil, false
}

func (mc *MaestroCache) Update(id string) error {
	e, ok := mc.entries.Get(id)
	if !ok {
		return fmt.Errorf("entry with id '%s' does not exist", id)
	}

	reg := e.(*entity.Registrant)
	reg.LastPing = time.Now()

	mc.entries.Set(id, reg, cache.DefaultExpiration)
	return nil
}

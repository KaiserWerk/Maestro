package cache

import (
	"testing"

	"github.com/KaiserWerk/Maestro/internal/entity"
)

type TestCache struct{}

var _ Cacher = &TestCache{}

func (mc *TestCache) Register(id, address string) error {
	return nil
}

func (mc *TestCache) Deregister(id string) bool {
	return true
}

func (mc *TestCache) Get(id string) (*entity.Registrant, bool) {
	return nil, true

}

func (mc *TestCache) Update(id string) error {
	return nil
}

func TestCacheImplementation(t *testing.T) {
	const (
		id      = "some-id"
		address = "http://some-address"
	)
	c := &TestCache{}

	if err := c.Register(id, address); err != nil {
		t.Fatalf("could not register: %s", err.Error())
	}

	if ok := c.Deregister(id); !ok {
		t.Fatalf("could not deregister")
	}

	if _, ok := c.Get(id); !ok {
		t.Fatalf("could not get entry")
	}

	if err := c.Update(id); err != nil {
		t.Fatalf("could not update: %s", err.Error())
	}
}

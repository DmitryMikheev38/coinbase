package cache

import (
	"coinbase/internal/models"
	"sync"
)

type Cache struct {
	tickStore sync.Map
}

func NewCache() *Cache {
	return &Cache{
		tickStore: sync.Map{},
	}
}

func (c *Cache) WriteTick(key string, value *models.Tick) {
	c.tickStore.Store(key, value)
}

func (c *Cache) ReadAndDeleteAllTicks() []*models.Tick {
	var res []*models.Tick
	c.tickStore.Range(func(key, value any) bool {
		tick := value.(*models.Tick)
		res = append(res, tick)
		c.tickStore.Delete(key)
		return true
	})
	return res
}

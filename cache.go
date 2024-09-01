package wbL0

import (
	"sync"
	"wbL0/tables"
)

type Cache struct {
	mu   sync.RWMutex
	data map[string]tables.Order
}

func NewCache() *Cache {
	return &Cache{data: make(map[string]tables.Order)}
}

func (c *Cache) SetOrder(order tables.Order) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.data[order.OrderUID] = order
}

func (c *Cache) GetOrderByUID(uid string) (tables.Order, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	order, ok := c.data[uid]
	return order, ok
}

func (c *Cache) RestoreFromDB(orders []tables.Order) {
	for _, order := range orders {
		c.data[order.OrderUID] = order
	}
}

func (c *Cache) CleanCache() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.data = make(map[string]tables.Order)
}

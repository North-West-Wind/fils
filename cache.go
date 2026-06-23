package main

type CacheLine struct {
	url string // Actual data
	ref bool   // Ref bit for CLOCK algorithm
}

type Cache struct {
	capacity   int                   // Max cache line stored
	mapping    map[uint64]*CacheLine // Mapping int ID to link
	order      []uint64              // Order of insertion
	orderIndex map[uint64]int        // Lookup table for key's index
	ptr        int                   // Pointer for CLOCK
}

func NewCache(capacity int) *Cache {
	cache := Cache{
		capacity:   capacity,
		mapping:    make(map[uint64]*CacheLine),
		order:      make([]uint64, 0),
		orderIndex: make(map[uint64]int),
		ptr:        0,
	}
	return &cache
}

func (c *Cache) Set(key uint64, value string) {
	// Remove from order if exists
	if c.Has(key) {
		index := c.orderIndex[key]
		c.order = append(c.order[:index], c.order[index+1:]...)
	}
	// Replace cache if needed
	if len(c.order) >= c.capacity {
		// Iterate starting from ptr
		for _, key := range append(c.order[c.ptr:], c.order[:c.ptr]...) {
			cacheLine := c.mapping[key]
			if cacheLine.ref {
				cacheLine.ref = false
			} else {
				c.Delete(key)
				break
			}
		}
	}
	// Insert value
	c.orderIndex[key] = len(c.order)
	c.order = append(c.order, key)
	c.mapping[key] = &CacheLine{
		url: value,
		ref: true,
	}
}

func (c *Cache) Delete(key uint64) {
	if _, exists := c.mapping[key]; exists {
		index := c.orderIndex[key]
		c.order = append(c.order[:index], c.order[index+1:]...)
		delete(c.orderIndex, key)
		delete(c.mapping, key)
		if c.ptr >= len(c.order) {
			c.ptr %= len(c.order)
		}
	}
}

func (c *Cache) Get(key uint64) *CacheLine {
	cacheLine, exists := c.mapping[key]
	if exists {
		cacheLine.ref = true
		return cacheLine
	} else {
		return nil
	}
}

func (c *Cache) Has(key uint64) bool {
	_, exists := c.mapping[key]
	return exists
}

package main

import "fmt"

type Cache struct {
	_cache map[string]string
}

func (c *Cache) initCache() {
	c._cache = make(map[string]string)

	fmt.Println("The cache created!")
}

func (c *Cache) set(key, value string) {
	c._cache[key] = value
}

func (c *Cache) get(key string) (string, bool) {
	value, ok := c._cache[key]
	return value, ok
}

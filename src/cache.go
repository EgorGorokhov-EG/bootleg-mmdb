package main

import (
	"log"
	"os"
)

type Cache struct {
	_cache map[string]string
	logger *log.Logger
}

func (c *Cache) init() {
	c._cache = make(map[string]string)
	c.logger = log.New(os.Stdout, "cache-logger", log.LstdFlags)

	c.logger.Println("The cache created!")
}

func (c *Cache) set(key, value string) {
	c._cache[key] = value
}

func (c *Cache) get(key string) (string, bool) {
	value, ok := c._cache[key]
	return value, ok
}

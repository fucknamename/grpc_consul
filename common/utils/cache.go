package utils

import (
	"context"
	"time"

	"github.com/allegro/bigcache/v3"
)

var c *bigcache.BigCache

func init() {
	c, _ = bigcache.New(context.Background(), bigcache.DefaultConfig(3*time.Minute))
}

func GetCache(key string) ([]byte, bool) {
	data, err := c.Get(key)
	return data, err == nil
}

func SetCache(key string, val []byte) {
	_ = c.Set(key, val)
}

func DelCache(key string) {
	_ = c.Delete(key)
}

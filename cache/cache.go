package cache

import (
	"github.com/dgraph-io/ristretto"
)

var (
	Cache *ristretto.Cache
)

func init() {
	var err error
	Cache, err = ristretto.NewCache(&ristretto.Config{
		NumCounters: 1000,
		MaxCost:     100,
		BufferItems: 64,
	})
	if err != nil {
		panic(err)
	}
}

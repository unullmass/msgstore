package cache

import (
	"github.com/dgraph-io/ristretto"
)

var (
	ReadCache *ristretto.Cache
)

func init() {
	var err error
	ReadCache, err = ristretto.NewCache(&ristretto.Config{
		NumCounters: 1000,
		MaxCost:     100,
		BufferItems: 64,
	})
	if err != nil {
		panic(err)
	}
}

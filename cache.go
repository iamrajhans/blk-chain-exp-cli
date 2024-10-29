// cache.go
package main

import (
	lru "github.com/hashicorp/golang-lru"
)

var blockCache *lru.Cache

func initCache() {
	var err error
	blockCache, err = lru.New(100) // Cache up to 100 items
	if err != nil {
		panic(err)
	}
}

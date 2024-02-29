package main

import (
	"runtime/debug"

	"github.com/coocood/freecache"
	"github.com/rs/zerolog/log"
)

func main() {

	cacheSize := 100 * 1024 * 1024
	cache := freecache.NewCache(cacheSize)
	debug.SetGCPercent(20)
	key := []byte("abc")
	val := []byte("def")
	expire := 60 // expire in 60 seconds
	cache.Set(key, val, expire)
	got, err := cache.Get(key)
	if err != nil {
		log.Info().Msgf("error is: %v", err)
	} else {
		log.Info().Msgf("%s", got)
	}
	affected := cache.Del(key)

	log.Info().Msgf("deleted key is: %v", affected)
	log.Info().Msgf("entry count is: %v", cache.EntryCount())
}

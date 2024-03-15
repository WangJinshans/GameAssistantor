package main

import (
	"time"

	"github.com/rs/zerolog/log"
)

func main() {

	c := make(chan struct{}, 1)

	go func(c <-chan struct{}) {
		<-c
		log.Info().Msg("notify go 1")
	}(c)

	go func(c <-chan struct{}) {
		<-c
		log.Info().Msg("notify go 2")
	}(c)

	go func(c <-chan struct{}) {
		<-c
		log.Info().Msg("notify go 3")
	}(c)

	close(c)
	time.Sleep(time.Second * 10)
}

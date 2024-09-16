package main

import (
	"github.com/117503445/goutils"
	"github.com/rs/zerolog/log"
)

func main() {
	goutils.InitZeroLog()

	log.Debug().Msg("debug")
	_, err := goutils.Exec("docker build -t 117503445/eventize-dev .")
	if err != nil {
		log.Fatal().Err(err).Msg("docker build failed")
	}

	
}

package main

import (
	"context"
	"net/http"

	"github.com/117503445/eventize/src/be/internal/rpc"
	"github.com/117503445/goutils"
	"github.com/rs/zerolog/log"
)

func main() {
	goutils.InitZeroLog()
	log.Debug().Msg("Hello, World!")

	client := rpc.NewHaberdasherProtobufClient("http://localhost:9090", &http.Client{})

	hat, err := client.MakeHat(context.Background(), &rpc.Size{Inches: -12})
	if err != nil {
		log.Fatal().Err(err).Msg("failed to make hat")
	}
	log.Info().Msgf("I have a nice new hat: %+v", hat)
}

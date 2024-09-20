package main

import (
	"context"
	"net/http"

	"github.com/117503445/eventize/src/be/internal/rpc"
	"github.com/117503445/goutils"
	"github.com/rs/zerolog/log"
	"github.com/twitchtv/twirp"
)

func main() {
	goutils.InitZeroLog()
	log.Debug().Msg("Hello, World!")

	client := rpc.NewEventizeProtobufClient("http://localhost:9090", &http.Client{}, twirp.WithClientPathPrefix("/rpc"))

	hat, err := client.MakeHat(context.Background(), &rpc.Size{Inches: 12})
	if err != nil {
		log.Fatal().Err(err).Msg("failed to make hat")
	}
	log.Info().Msgf("I have a nice new hat: %+v", hat)

	resp, err := client.CreateEvent(context.Background(), &rpc.CreateEventRequest{
		Type: "test",
	})
	if err != nil {
		log.Fatal().Err(err).Msg("failed to create event")
	}
	log.Info().Msgf("Event created: %+v", resp)
}

package main

import (
	"context"
	"net/http"

	"github.com/117503445/eventize/src/be/internal/common"
	"github.com/117503445/eventize/src/be/internal/rpc"
	"github.com/117503445/goutils"
	"github.com/rs/zerolog/log"
	"github.com/twitchtv/twirp"
	"google.golang.org/protobuf/types/known/emptypb"
)

func main() {
	goutils.InitZeroLog(goutils.WithProduction{
		FileName: "agent",
	})
	var err error
	log.Debug().Msg("Hello, World!")

	client := rpc.NewEventizeProtobufClient("http://server:9090", &http.Client{}, twirp.WithClientPathPrefix("/rpc"))

	var buildInfo *rpc.BuildInfo
	err = common.NewRetry().Execute(func() error {
		var err error
		buildInfo, err = client.GetBuildInfo(context.Background(), &emptypb.Empty{})
		return err
	})
	if err != nil {
		log.Fatal().Err(err).Msg("failed to get build info")
	}
	log.Info().Interface("buildInfo", buildInfo).Msg("Got build info")

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

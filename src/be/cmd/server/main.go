package main

import (
	"net/http"

	"github.com/117503445/eventize/src/be/internal/rpc"
	"github.com/117503445/eventize/src/be/internal/server"
	"github.com/117503445/goutils"
	"github.com/rs/zerolog/log"
)

func main() {
	goutils.InitZeroLog()
	log.Debug().Msg("Hello, World!")

	server := &server.Server{} // implements Haberdasher interface
	twirpHandler := rpc.NewHaberdasherServer(server)

	if err := http.ListenAndServe(":9090", twirpHandler); err != nil {
		panic(err)
	}
}

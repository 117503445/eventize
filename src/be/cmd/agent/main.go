package main

import (
	"bytes"
	"context"
	"errors"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/117503445/eventize/src/be/internal/agent"
	"github.com/117503445/eventize/src/be/internal/common"
	"github.com/117503445/eventize/src/be/internal/rpc"
	"github.com/117503445/goutils"
	"github.com/coder/websocket"
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

	go func() {
		rpcServer := &agent.Server{}
		twirpHandler := rpc.NewEventizeAgentServer(rpcServer)
		log.Info().Msg("Starting server")
		err := http.ListenAndServe(":9090", twirpHandler)
		if err != nil {
			log.Fatal().Err(err).Msg("failed to start server")
		}
	}()

	client := rpc.NewEventizeServerProtobufClient("http://server:9090", &http.Client{}, twirp.WithClientPathPrefix("/rpc"))

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

	go ws()

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

	select {}
}

// create TCP connection to server by HTTP request
func createTcpConnByHTTPRequest(req []byte) (net.Conn, error) {
	// get first line of HTTP request
	// line := lines[:bytes.IndexByte(req, '\n')]
	line := string(req[:bytes.IndexByte(req[:], '\n')])

	// split the first line by spaces
	parts := strings.Split(line, " ")
	if len(parts) < 2 {
		return nil, errors.New("invalid HTTP request line")
	}

	// extract the host from the URL
	url := parts[1]
	host := strings.Split(url, "/")[2]

	log.Debug().Discard().Str("address", host).Msg("connecting to server")

	// create TCP connection
	conn, err := net.Dial("tcp", host)
	if err != nil {
		return nil, err
	}

	return conn, nil
}

func HandleHttpProxyRequest(c *websocket.Conn) error {
	_, req, err := c.Read(context.Background())
	if err != nil {
		log.Error().Err(err).Msg("Error reading request")
		return err
	}

	// fallback try to send 500 error to websocket connection
	fallback := func() error {
		err := c.Write(context.Background(), websocket.MessageBinary, []byte("HTTP/1.1 500 Internal Server Error\r\n\r\n"))
		if err != nil {
			log.Error().Err(err).Msg("Error writing to client")
			return err
		} else {
			return nil
		}
	}

	serverConn, err := createTcpConnByHTTPRequest(req)
	if err != nil {
		log.Error().Err(err).Msg("Error creating connection to server")
		if err := fallback(); err != nil {
			return err
		}
	}
	defer serverConn.Close()

	_, err = serverConn.Write(req)
	if err != nil {
		log.Error().Err(err).Msg("Error writing to server")
		if err := fallback(); err != nil {
			return err
		}
	}
	resp, err := common.ReadHttpFromTcp(serverConn)
	if err != nil {
		log.Error().Err(err).Msg("Error reading from server")
		if err := fallback(); err != nil {
			return err
		}
	}

	err = c.Write(context.Background(), websocket.MessageBinary, resp)
	if err != nil {
		log.Error().Err(err).Msg("Error writing to client")
		return err
	}
	return nil
}

func ws() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	id := goutils.UUID4()

	c, _, err := websocket.Dial(ctx, "ws://server:9090/ws", &websocket.DialOptions{
		HTTPHeader: map[string][]string{"X-Id": {id}},
	})
	if err != nil {
		log.Fatal().Err(err).Msg("failed to dial")
	}

	log.Info().Str("Subprotocol", c.Subprotocol()).Msg("ws connected")

	for {
		if err := HandleHttpProxyRequest(c); err != nil {
			break
		}
	}

	log.Info().Msg("ws disconnected")
}

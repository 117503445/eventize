package main

import (
	// "net"
	"net/http"
	"time"

	// "github.com/117503445/eventize/src/be/internal/rpc"
	// "github.com/117503445/eventize/src/be/internal/server"
	"github.com/117503445/goutils"
	"github.com/rs/zerolog/log"

	"context"
	"fmt"
	"io"

	"golang.org/x/time/rate"

	"github.com/coder/websocket"
)

func main() {
	goutils.InitZeroLog()
	log.Debug().Msg("Hello, World!")

	// server := &server.Server{} // implements Haberdasher interface
	// twirpHandler := rpc.NewHaberdasherServer(server)

	// if err := http.ListenAndServe(":9090", twirpHandler); err != nil {
	// 	panic(err)
	// }

	s := &http.Server{
		Handler: echoServer{
			logf: log.Printf,
		},
		ReadTimeout:  time.Second * 60,
		WriteTimeout: time.Second * 60,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		s.Handler.ServeHTTP(w, r)
	})

	muxServer := &http.Server{
		Addr:         ":9090",
		Handler:      mux,
	}
	if err := muxServer.ListenAndServe(); err != nil {
		log.Fatal().Err(err).Msg("failed to serve")
	}
	// l, err := net.Listen("tcp", ":9090")
	// if err != nil {
	// 	log.Fatal().Err(err).Msg("failed to listen")
	// }


	// err = s.Serve(l)
	// if err != nil {
	// 	log.Fatal().Err(err).Msg("failed to serve")
	// }

}

// echoServer is the WebSocket echo server implementation.
// It ensures the client speaks the echo subprotocol and
// only allows one message every 100ms with a 10 message burst.
type echoServer struct {
	// logf controls where logs are sent.
	logf func(f string, v ...interface{})
}

func (s echoServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c, err := websocket.Accept(w, r, &websocket.AcceptOptions{
		// Subprotocols: []string{"echo"},
		InsecureSkipVerify: true,
	})
	if err != nil {
		s.logf("%v", err)
		return
	}
	defer c.CloseNow()

	// if c.Subprotocol() != "echo" {
	// 	c.Close(websocket.StatusPolicyViolation, "client must speak the echo subprotocol")
	// 	return
	// }

	log.Debug().Msg("accepted connection")

	l := rate.NewLimiter(rate.Every(time.Millisecond*100), 10)
	for {
		err = echo(c, l)
		if websocket.CloseStatus(err) == websocket.StatusNormalClosure {
			return
		}
		if err != nil {
			log.Error().Err(err).Str("remote", r.RemoteAddr).Msg("failed to echo")
			return
		}
	}
}

// echo reads from the WebSocket connection and then writes
// the received message back to it.
// The entire function has 10s to complete.
func echo(c *websocket.Conn, l *rate.Limiter) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	err := l.Wait(ctx)
	if err != nil {
		return err
	}

	log.Debug().Msg("reading message")

	typ, r, err := c.Reader(ctx)
	if err != nil {
		return err
	}

	log.Debug().Msg("writing message")

	w, err := c.Writer(ctx, typ)
	if err != nil {
		return err
	}

	log.Debug().Msg("copying message")

	_, err = io.Copy(w, r)
	if err != nil {
		return fmt.Errorf("failed to io.Copy: %w", err)
	}

	err = w.Close()
	return err
}

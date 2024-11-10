package main

import (
	"context"
	"embed"
	"io/fs"
	"net"
	"net/http"
	"os"

	"github.com/117503445/goutils"
	"github.com/coder/websocket"
	"github.com/rs/zerolog/log"
	"github.com/twitchtv/twirp"
	"golang.org/x/net/webdav"

	"github.com/117503445/eventize/src/be/internal/common"
	"github.com/117503445/eventize/src/be/internal/rpc"
	"github.com/117503445/eventize/src/be/internal/server"
)

//go:embed all:dist
var staticFiles embed.FS

func main() {
	goutils.InitZeroLog(goutils.WithProduction{
		FileName: "server",
	})
	log.Debug().Interface("buildInfo", common.GetBuildInfo()).Msg("Eventize server")

	rpcServer := server.NewServer()
	twirpHandler := rpc.NewEventizeServerServer(rpcServer, twirp.WithServerPathPrefix("/rpc"))

	log.Debug().Str("prefix", twirpHandler.PathPrefix()).Msg("twirp handler path prefix")

	rr, err := staticFiles.ReadDir("dist/assets")
	if err != nil {
		log.Fatal().Err(err).Msg("failed to read assets")
	}
	for _, r := range rr {
		log.Debug().Discard().Str("name", r.Name()).Msg("static file")
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/ws", wsHandler)
	mux.HandleFunc(twirpHandler.PathPrefix(), twirpHandler.ServeHTTP)

	dirWebdav := "./data/webdav"
	if err := os.MkdirAll(dirWebdav, os.ModePerm); err != nil {
		log.Fatal().Err(err).Msg("failed to create webdav dir")
	}

	webdavHandler := &webdav.Handler{
		FileSystem: webdav.Dir("./data/webdav"),
		LockSystem: webdav.NewMemLS(),
		Prefix:     "/webdav",
	}
	mux.HandleFunc("/webdav/", func(w http.ResponseWriter, r *http.Request) {
		log.Debug().Str("method", r.Method).Str("path", r.URL.Path).Msg("webdav")
		webdavHandler.ServeHTTP(w, r)
	})

	feFs, err := fs.Sub(staticFiles, "dist")
	if err != nil {
		log.Fatal().Err(err).Msg("failed to create assets fs")
	}
	mux.Handle("/", http.StripPrefix("/", http.FileServer(http.FS(feFs))))

	// CORS 中间件
	enableCORS := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

			// 如果是预检请求，直接返回
			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusOK)
				return
			}

			next.ServeHTTP(w, r)
		})
	}

	muxServer := &http.Server{
		Addr:    ":9090",
		Handler: enableCORS(mux),
	}
	if err := muxServer.ListenAndServe(); err != nil {
		log.Fatal().Err(err).Msg("failed to serve")
	}
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	id := r.Header.Get("X-Id")
	log.Debug().Str("id", id).Str("RemoteAddr", r.RemoteAddr).
		Msg("received connection")

	c, err := websocket.Accept(w, r, &websocket.AcceptOptions{
		InsecureSkipVerify: true,
	})
	if err != nil {
		log.Error().Err(err).Msg("failed to accept")
		return
	}
	log.Info().Str("Subprotocol", c.Subprotocol()).Msg("ws connected")

	go func() {
		listener, err := net.Listen("tcp", ":8081")
		if err != nil {
			log.Fatal().Err(err).Msg("Error creating listener")
		}
		defer listener.Close()
		log.Info().Msg("HTTP Proxy Listening on :8081")
		for {
			clientConn, err := listener.Accept()
			if err != nil {
				log.Error().Err(err).Msg("Error accepting connection")
				continue
			}
			err = HandleHttpProxyRequest(clientConn, c)
			if err != nil {
				break
			}
		}
		log.Info().Err(err).Msg("Closing HTTP Proxy connection")
	}()
}

func HandleHttpProxyRequest(clientConn net.Conn, c *websocket.Conn) error {
	defer clientConn.Close()

	log.Debug().Discard().Msg("Accepted HTTP Proxy connection")

	req, err := common.ReadHttpFromTcp(clientConn)
	if err != nil {
		log.Error().Err(err).Msg("Error reading Client Request")
		return nil
	}

	err = c.Write(context.Background(), websocket.MessageBinary, req)
	if err != nil {
		log.Error().Err(err).Msg("Error writing to WebSocket Connection")
		return err
	}

	_, resp, err := c.Read(context.Background())
	if err != nil {
		log.Error().Err(err).Msg("Error reading from WebSocket Connection")
		return err
	}

	_, err = clientConn.Write(resp)
	if err != nil {
		log.Error().Err(err).Msg("Error writing to Client Connection")
		return nil
	}

	log.Debug().Discard().Msg("Closed HTTP Proxy connection")
	return nil
}

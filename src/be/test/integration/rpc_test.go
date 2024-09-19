package integration_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/117503445/eventize/src/be/internal/rpc"
	"github.com/117503445/goutils"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
	"github.com/twitchtv/twirp"
	"google.golang.org/protobuf/types/known/emptypb"
)

func TestMain(m *testing.M) {
	goutils.InitZeroLog(goutils.WithNoColor{})
	log.Debug().Msg("Hello, World!")
	m.Run()
}

func TestGetServerMeta(t *testing.T) {
	var err error
	ast := assert.New(t)
	client := rpc.NewEventizeProtobufClient("http://localhost:9090", &http.Client{}, twirp.WithClientPathPrefix("/rpc"))
	meta, err := client.GetBuildInfo(context.Background(), &emptypb.Empty{})
	ast.Nil(err)

	ast.NotEmpty(meta.Version)
}

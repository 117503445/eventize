package agent

import (
	"context"

	"github.com/117503445/eventize/src/be/internal/common"
	"github.com/117503445/eventize/src/be/internal/rpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

type Server struct {
	rpc.EventizeAgent
}

func (s *Server) GetBuildInfo(context.Context, *emptypb.Empty) (meta *rpc.BuildInfo, err error) {
	return &rpc.BuildInfo{
		Version: common.GetBuildInfo().Version,
	}, nil
}

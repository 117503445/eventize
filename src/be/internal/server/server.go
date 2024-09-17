package server

import (
	"context"
	"math/rand"

	"github.com/117503445/eventize/src/be/internal/rpc"
	"github.com/twitchtv/twirp"
)

// Server implements the Haberdasher service
type Server struct{}

func (s *Server) MakeHat(ctx context.Context, size *rpc.Size) (hat *rpc.Hat, err error) {
	if size.Inches <= 0 {
		return nil, twirp.InvalidArgumentError("inches", "I can't make a hat that small!")
	}
	return &rpc.Hat{
		Inches: size.Inches,
		Color:  []string{"white", "black", "brown", "red", "blue"}[rand.Intn(5)],
		Name:   []string{"bowler", "baseball cap", "top hat", "derby"}[rand.Intn(4)],
	}, nil
}
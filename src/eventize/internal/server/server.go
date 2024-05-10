package server

import (
	"context"
	"fmt"
	"net"

	"github.com/117503445/eventize/src/eventize/internal/pb"
	"google.golang.org/grpc"
)

type Server struct {
	pb.UnimplementedGreeterServer

	grpcServer *grpc.Server
	port       int
}

func NewServer(port int) *Server {
	s := grpc.NewServer()
	pb.RegisterGreeterServer(s, &Server{})
	return &Server{grpcServer: s, port: port}
}

func (s *Server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	return &pb.HelloReply{Message: "Hello " + in.Name}, nil
}

func (s *Server) Serve() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", s.port))
	if err != nil {
		panic(err)
	}


	err = s.grpcServer.Serve(lis)
	if err != nil {
		panic(err)
	}

}

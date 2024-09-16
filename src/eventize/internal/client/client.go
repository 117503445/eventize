package client

import (
	"github.com/117503445/eventize/src/eventize/internal/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	// "google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	pb.GreeterClient

	conn *grpc.ClientConn
}

func NewClient(addr string) *Client {
	// conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))

	// https
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(
		credentials.NewClientTLSFromCert(nil, "")),
	)
	if err != nil {
		panic(err)
	}

	grpcClient := pb.NewGreeterClient(conn)
	return &Client{grpcClient, conn}
}

func (c *Client) Close() {
	c.conn.Close()
}

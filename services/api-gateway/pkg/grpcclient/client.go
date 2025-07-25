package grpcclient

import (
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	_defaultAddr = ":80"
)

type Client struct {
	Conn *grpc.ClientConn
	notify chan error
}

func New(target string, opts ...grpc.DialOption) (*Client, error) {
	if len(opts) == 0 {
		opts = append(opts, grpc.WithTransportCredentials((insecure.NewCredentials())))
	}

	// opts = append(opts, grpc.WithBlock())

	conn, err := grpc.NewClient(target, opts...)
	if err != nil {
		return nil, fmt.Errorf("grpcclient - New - failed to connect to %s: %w", target, err)
	}

	return &Client{
		Conn: conn,
	}, nil
}

func (c *Client) Close() error {
	if c.Conn != nil {
		return c.Conn.Close()
	}
	return nil
}

package client

import (
	"context"
	"fmt"

	"github.com/angelini/gotemplate/pkg/pb"
	"google.golang.org/grpc"
)

type Client struct {
	conn    *grpc.ClientConn
	example pb.ExampleClient
}

func NewClient(ctx context.Context, server string) (*Client, error) {
	conn, err := grpc.DialContext(ctx, server, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		return nil, err
	}

	return &Client{conn: conn, example: pb.NewExampleClient(conn)}, nil
}

func (c *Client) Close() {
	c.conn.Close()
}

type Output struct {
	Static int64
	FromDb int64
}

func (c *Client) GetAll(ctx context.Context) (Output, error) {
	output := Output{}

	result, err := c.example.Static(ctx, &pb.ExampleRequest{})
	if err != nil {
		return output, fmt.Errorf("could not get static: %s", err)
	}
	output.Static = result.GetData()

	result, err = c.example.FromDb(ctx, &pb.ExampleRequest{})
	if err != nil {
		return output, fmt.Errorf("could not get from db: %s", err)
	}
	output.FromDb = result.GetData()

	return output, nil
}

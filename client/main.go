package main

import (
	"context"
	"time"

	pb "github.com/angelini/gotemplate/proto"

	"go.uber.org/zap"
	"google.golang.org/grpc"
)

const (
	address = "localhost:5051"
)

func main() {
	logger, _ := zap.NewDevelopment()
	defer logger.Sync()

	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		logger.Fatal("could not connect to server", zap.String("address", address))
	}
	defer conn.Close()

	c := pb.NewExampleClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	r, err := c.Test(ctx, &pb.TestRequest{})
	if err != nil {
		logger.Fatal("could not get status", zap.Error(err))
	}
	logger.Info("data", zap.Int64("value", r.GetData()))
}

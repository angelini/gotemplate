package main

import (
	"context"
	"os"
	"time"

	pb "github.com/angelini/gotemplate/proto"

	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func main() {
	logger, _ := zap.NewDevelopment()
	defer logger.Sync()

	server := os.Getenv("SERVER")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	conn, err := grpc.DialContext(ctx, server, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		logger.Fatal("could not connect to server", zap.String("address", server))
	}
	defer conn.Close()

	c := pb.NewExampleClient(conn)

	ctx, cancel = context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	r, err := c.Test(ctx, &pb.TestRequest{})
	if err != nil {
		logger.Fatal("could not get status", zap.Error(err))
	}
	logger.Info("data", zap.Int64("value", r.GetData()))
}

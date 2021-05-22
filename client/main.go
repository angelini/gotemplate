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

	c := pb.NewStatusClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	r, err := c.Status(ctx, &pb.StatusRequest{})
	if err != nil {
		logger.Fatal("could not get status", zap.Error(err))
	}
	logger.Info("status", zap.Bool("up", r.GetUp()))

	r, err = c.DbStatus(ctx, &pb.StatusRequest{})
	if err != nil {
		logger.Fatal("could not get db-status", zap.Error(err))
	}
	logger.Info("db-status", zap.Bool("up", r.GetUp()))
}

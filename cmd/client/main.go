package main

import (
	"context"
	"os"
	"time"

	"github.com/angelini/gotemplate/pkg/client"

	"go.uber.org/zap"
)

func main() {
	log, _ := zap.NewDevelopment()
	defer log.Sync()

	server := os.Getenv("SERVER")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	c, err := client.NewClient(ctx, server)
	if err != nil {
		log.Fatal("could not connect to server", zap.String("server", server))
	}
	defer c.Close()

	output, err := c.GetAll(ctx)
	if err != nil {
		log.Fatal("could not fetch data", zap.Error(err))
	}

	log.Info("static", zap.Int64("data", output.Static))
	log.Info("from DB", zap.Int64("data", output.FromDb))
}

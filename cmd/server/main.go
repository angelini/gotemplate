package main

import (
	"context"
	"net"
	"os"

	"go.uber.org/zap"

	"github.com/angelini/gotemplate/pkg/api"
	"github.com/angelini/gotemplate/pkg/server"
	"github.com/jackc/pgx/v4/pgxpool"
)

func main() {
	ctx := context.Background()

	log, _ := zap.NewDevelopment()
	defer log.Sync()

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("missing PORT env variable")
	}

	dbUri := os.Getenv("DB_URI")

	listen, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatal("failed to listen", zap.String("protocol", "tcp"), zap.String("port", port))
	}

	pool, err := pgxpool.Connect(ctx, dbUri)
	if err != nil {
		log.Fatal("cannot connect to DB", zap.String("uri", dbUri))
	}
	defer pool.Close()

	s := server.NewServer(log)

	log.Info("register Example")
	example := &api.Example{
		Log:  log,
		Pool: pool,
	}
	s.RegisterExample(ctx, example)

	log.Info("start server", zap.String("port", port))
	if err := s.Serve(listen); err != nil {
		log.Fatal("failed to serve", zap.Error(err))
	}
}

package main

import (
	"context"
	"net"
	"os"

	pb "github.com/angelini/gotemplate/proto"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"github.com/jackc/pgx/v4/pgxpool"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type MissingConfigError struct{}

func (e *MissingConfigError) Error() string {
	return "missing config in the request context"
}

type server struct {
	pb.UnimplementedStatusServer
	log  *zap.Logger
	pool *pgxpool.Pool
}

func (s *server) Status(ctx context.Context, in *pb.StatusRequest) (*pb.StatusResponse, error) {
	return &pb.StatusResponse{Up: true}, nil
}

func (s *server) DbStatus(ctx context.Context, in *pb.StatusRequest) (*pb.StatusResponse, error) {
	conn, err := s.pool.Acquire(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	return &pb.StatusResponse{Up: true}, nil
}

func main() {
	log, _ := zap.NewDevelopment()
	defer log.Sync()

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("missing PORT env variable")
	}

	listen, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatal("failed to listen", zap.String("protocol", "tcp"), zap.String("port", port))
	}

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(
			grpc_middleware.ChainUnaryServer(
				grpc_zap.UnaryServerInterceptor(log),
				grpc_recovery.UnaryServerInterceptor(),
			),
		),
	)

	dbUri := os.Getenv("DB_URI")
	pool, err := pgxpool.Connect(context.Background(), dbUri)
	if err != nil {
		log.Fatal("cannot connect to DB", zap.String("uri", dbUri))
	}
	defer pool.Close()

	log.Info("registering StatusServer", zap.String("port", port))
	pb.RegisterStatusServer(grpcServer, &server{
		log:  log,
		pool: pool,
	})
	if err := grpcServer.Serve(listen); err != nil {
		log.Fatal("failed to server StatusServer", zap.Error(err))
	}
}

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
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

type MissingConfigError struct{}

func (e *MissingConfigError) Error() string {
	return "missing config in the request context"
}

type appServer struct {
	pb.UnimplementedExampleServer
	log  *zap.Logger
	pool *pgxpool.Pool
}

func (s *appServer) Test(ctx context.Context, in *pb.TestRequest) (*pb.TestResponse, error) {
	return &pb.TestResponse{Data: 42}, nil
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

	log.Info("register HealthServer")
	healthServer := health.NewServer()
	healthpb.RegisterHealthServer(grpcServer, healthServer)

	log.Info("register StatusServer")
	pb.RegisterExampleServer(grpcServer, &appServer{
		log:  log,
		pool: pool,
	})
	healthServer.SetServingStatus("gotemplate.server.Example", healthpb.HealthCheckResponse_SERVING)

	log.Info("start server", zap.String("port", port))
	if err := grpcServer.Serve(listen); err != nil {
		log.Fatal("failed to server StatusServer", zap.Error(err))
	}
}

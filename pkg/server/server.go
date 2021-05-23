package server

import (
	"context"
	"net"

	"github.com/angelini/gotemplate/pkg/api"
	"github.com/angelini/gotemplate/pkg/pb"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

type Server struct {
	Grpc   *grpc.Server
	Health *health.Server
}

func NewServer(log *zap.Logger) *Server {
	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(
			grpc_middleware.ChainUnaryServer(
				grpc_zap.UnaryServerInterceptor(log),
				grpc_recovery.UnaryServerInterceptor(),
			),
		),
	)

	log.Info("register HealthServer")
	healthServer := health.NewServer()
	healthpb.RegisterHealthServer(grpcServer, healthServer)

	return &Server{Grpc: grpcServer, Health: healthServer}
}

func (s *Server) RegisterExample(ctx context.Context, example *api.Example) {
	pb.RegisterExampleServer(s.Grpc, example)
	example.HealthMonitor(ctx, s.Health)
}

func (s *Server) Serve(lis net.Listener) error {
	return s.Grpc.Serve(lis)
}

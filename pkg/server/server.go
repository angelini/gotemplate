package server

import (
	"context"
	"net"
	"time"

	"github.com/angelini/gotemplate/pkg/api"
	"github.com/angelini/gotemplate/pkg/pb"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"github.com/jackc/pgx/v4/pgxpool"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

type Server struct {
	log    *zap.Logger
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

	return &Server{log: log, Grpc: grpcServer, Health: healthServer}
}

func (s *Server) MonitorDbPool(ctx context.Context, pool *pgxpool.Pool) {
	ticker := time.NewTicker(time.Second)

	go func() {
		for {
			select {
			case <-ctx.Done():
				s.Health.SetServingStatus("gotemplate.server.Example", healthpb.HealthCheckResponse_NOT_SERVING)
			case <-ticker.C:
				ctxTimeout, cancel := context.WithTimeout(ctx, 800*time.Millisecond)

				status := healthpb.HealthCheckResponse_SERVING
				err := pool.Ping(ctxTimeout)
				if err != nil {
					status = healthpb.HealthCheckResponse_NOT_SERVING
				}
				cancel()

				s.Health.SetServingStatus("gotemplate.server.Example", status)
			}
		}
	}()
}

func (s *Server) RegisterExample(ctx context.Context, example *api.Example) {
	pb.RegisterExampleServer(s.Grpc, example)
}

func (s *Server) Serve(lis net.Listener) error {
	return s.Grpc.Serve(lis)
}

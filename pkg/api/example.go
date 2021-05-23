package api

import (
	"context"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
	"go.uber.org/zap"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"

	pb "github.com/angelini/gotemplate/pkg/pb"
)

type Example struct {
	pb.UnimplementedExampleServer

	Log  *zap.Logger
	Pool *pgxpool.Pool
}

func (e *Example) Static(ctx context.Context, in *pb.ExampleRequest) (*pb.ExampleResponse, error) {
	return &pb.ExampleResponse{Data: 42}, nil
}

func (e *Example) FromDb(ctx context.Context, in *pb.ExampleRequest) (*pb.ExampleResponse, error) {
	conn, err := e.Pool.Acquire(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	rows, err := conn.Query(ctx, "SELECT count(*) FROM information_schema.schemata;")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	rows.Next()
	values, err := rows.Values()
	if err != nil {
		return nil, err
	}

	return &pb.ExampleResponse{Data: values[0].(int64)}, nil
}

func (e *Example) HealthMonitor(ctx context.Context, healthServer *health.Server) {
	ticker := time.NewTicker(time.Second)

	go func() {
		for {
			select {
			case <-ctx.Done():
				healthServer.SetServingStatus("gotemplate.server.Example", healthpb.HealthCheckResponse_NOT_SERVING)
			case <-ticker.C:
				status := healthpb.HealthCheckResponse_SERVING
				err := e.Pool.Ping(ctx)
				if err != nil {
					status = healthpb.HealthCheckResponse_NOT_SERVING
				}
				healthServer.SetServingStatus("gotemplate.server.Example", status)
			}
		}
	}()
}

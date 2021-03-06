package api

import (
	"context"

	"go.uber.org/zap"

	"github.com/angelini/gotemplate/pkg/pb"
)

type Example struct {
	pb.UnimplementedExampleServer

	Log    *zap.Logger
	DbConn DbConnector
}

func (e *Example) Static(ctx context.Context, in *pb.ExampleRequest) (*pb.ExampleResponse, error) {
	return &pb.ExampleResponse{Data: 42}, nil
}

func (e *Example) FromDb(ctx context.Context, in *pb.ExampleRequest) (*pb.ExampleResponse, error) {
	conn, cancel, err := e.DbConn.Connect(ctx)
	if err != nil {
		return nil, err
	}
	defer cancel()

	rows, err := conn.Query(ctx, "SELECT count(*) FROM example.t1;")
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

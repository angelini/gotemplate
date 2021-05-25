package test

import (
	"context"
	"os"
	"testing"

	"github.com/angelini/gotemplate/pkg/api"
	"github.com/angelini/gotemplate/pkg/pb"
	"github.com/jackc/pgx/v4"
	"go.uber.org/zap"
)

type DbTestConnector struct {
	dbUri string
}

func (d *DbTestConnector) Connect(ctx context.Context) (*pgx.Conn, func(), error) {
	conn, err := pgx.Connect(ctx, d.dbUri)
	if err != nil {
		return nil, nil, err
	}

	tx, err := conn.Begin(ctx)
	if err != nil {
		return nil, nil, err
	}

	return conn, func() { tx.Rollback(ctx) }, nil
}

var (
	log, _ = zap.NewDevelopment()
	dbConn = &DbTestConnector{dbUri: os.Getenv("DB_URI")}
)

func TestStatic(t *testing.T) {
	a := api.Example{
		Log:    log,
		DbConn: dbConn,
	}

	response, err := a.Static(context.Background(), &pb.ExampleRequest{})
	if err != nil {
		t.Errorf("%s", err)
	}

	data := response.GetData()
	if data != 42 {
		t.Errorf("incorrect data, got: %d, want: %d", data, 42)
	}
}

func TestFromDb(t *testing.T) {
	a := api.Example{
		Log:    log,
		DbConn: dbConn,
	}

	response, err := a.FromDb(context.Background(), &pb.ExampleRequest{})
	if err != nil {
		t.Errorf("%s", err)
	}

	data := response.GetData()
	if data != 6 {
		t.Errorf("incorrect data, got: %d, want: %d", data, 6)
	}
}

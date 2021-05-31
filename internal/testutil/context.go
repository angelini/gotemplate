package testutil

import (
	"context"
	"os"
	"testing"

	"github.com/angelini/gotemplate/pkg/api"
	"github.com/jackc/pgx/v4"
)

type TestCtx struct {
	t      *testing.T
	dbConn *DbTestConnector
	ctx    context.Context
}

func NewTestCtx(t *testing.T) TestCtx {
	ctx := context.Background()

	dbConn, err := newDbTestConnector(ctx, os.Getenv("DB_URI"))
	if err != nil {
		t.Errorf("connecting to DB: %w", err)
	}

	return TestCtx{
		t:      t,
		dbConn: dbConn,
		ctx:    ctx,
	}
}

func (tc *TestCtx) Connector() api.DbConnector {
	return tc.dbConn
}

func (tc *TestCtx) Connect() *pgx.Conn {
	conn, _, err := tc.dbConn.Connect(tc.ctx)
	if err != nil {
		tc.Errorf("connecting to db: %w", err)
	}
	return conn
}

func (tc *TestCtx) Errorf(format string, args ...interface{}) {
	tc.t.Errorf(format, args...)
}

func (tc *TestCtx) Context() context.Context {
	return tc.ctx
}

func (tc *TestCtx) Close() {
	tc.dbConn.close(tc.ctx)
}

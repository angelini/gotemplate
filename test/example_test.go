package test

import (
	"testing"

	utils "github.com/angelini/gotemplate/internal/testutil"
	"github.com/angelini/gotemplate/pkg/api"
	"github.com/angelini/gotemplate/pkg/pb"

	"go.uber.org/zap"
)

func writeT1(tc *utils.TestCtx, id int32, value string) {
	conn := tc.Connect()

	_, err := conn.Exec(tc.Context(), `
		INSERT INTO example.t1 (id, val)
		VALUES ($1, $2)
		`, id, value)

	if err != nil {
		tc.Errorf("error inserting into t1: %w", err)
	}
}

var (
	log, _ = zap.NewDevelopment()
)

func TestStatic(t *testing.T) {
	tc := utils.NewTestCtx(t)
	defer tc.Close()

	a := api.Example{
		Log:    log,
		DbConn: tc.Connector(),
	}

	response, err := a.Static(tc.Context(), &pb.ExampleRequest{})
	if err != nil {
		t.Errorf("%s", err)
	}

	data := response.GetData()
	if data != 42 {
		t.Errorf("incorrect data, got: %d, want: %d", data, 42)
	}
}

func TestFromDb(t *testing.T) {
	tc := utils.NewTestCtx(t)
	defer tc.Close()

	a := api.Example{
		Log:    log,
		DbConn: tc.Connector(),
	}

	writeT1(&tc, 1, "foo")
	writeT1(&tc, 2, "bar")

	response, err := a.FromDb(tc.Context(), &pb.ExampleRequest{})
	if err != nil {
		t.Errorf("%s", err)
	}

	data := response.GetData()
	if data != 2 {
		t.Errorf("incorrect data, got: %d, want: %d", data, 2)
	}
}

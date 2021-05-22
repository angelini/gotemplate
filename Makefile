.PHONY: install build server client health-probes

export PORT=:5051
export DB_URI=postgres://postgres@10.1.1.3:5432/postgres

install:
	go install google.golang.org/protobuf/cmd/protoc-gen-go
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc
	go install github.com/grpc-ecosystem/grpc-health-probe

proto/%.pb.go: proto/%.proto
	protoc --go_out=. --go_opt=paths=source_relative $^

proto/%_grpc.pb.go: proto/%.proto
	protoc --go-grpc_out=. --go-grpc_opt=paths=source_relative $^

bin/%: %/main.go
	go build -o $@ $^

build: proto/service.pb.go proto/service_grpc.pb.go bin/server bin/client

server:
	go run server/main.go

client:
	go run client/main.go

health-probes:
	grpc-health-probe -addr localhost:5051
	grpc-health-probe -addr localhost:5051 -service gotemplate.server.Example

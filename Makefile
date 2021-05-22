.PHONY: build server client

export PORT=:5051
export DB_URI=postgres://postgres@10.1.1.3:5432/postgres

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

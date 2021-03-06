DB_HOST := 10.1.1.3
DB_URI := postgres://postgres@$(DB_HOST):5432/example
export DB_URI

PKG_GO_FILES := $(shell find pkg/ -type f -name '*.go')
MIGRATE_DIR := ./migrations
SERVICE := gotemplate.server

.PHONY: install test build server client health
.PHONY: clear-k8s build-k8s deploy-k8s client-k8s health-k8s
.PHONY: migrate migrate-create

install:
	go install google.golang.org/protobuf/cmd/protoc-gen-go
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc
	go install github.com/grpc-ecosystem/grpc-health-probe
	go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate

test:
	cd test && go test

pkg/pb/%.pb.go: pkg/pb/%.proto
	protoc --go_out=. --go_opt=paths=source_relative $^

pkg/pb/%_grpc.pb.go: pkg/pb/%.proto
	protoc --go-grpc_out=. --go-grpc_opt=paths=source_relative $^

bin/%: cmd/%/main.go $(PKG_GO_FILES)
	go build -o $@ $<

build: pkg/pb/example.pb.go pkg/pb/example_grpc.pb.go bin/server bin/client

server: export PORT=:5051
server:
	go run cmd/server/main.go

client: export SERVER=localhost:5051
client:
	go run cmd/client/main.go

health: export SERVER=localhost:5051
health:
	grpc-health-probe -addr $(SERVER)
	grpc-health-probe -addr $(SERVER) -service $(SERVICE)

k8s-clear:
	kubectl delete --all service
	kubectl delete --all pod --grace-period 0 --force

k8s-build: build
	podman build -f Dockerfile -t "gotemplate:server"
	podman save -o /tmp/gotemplate_server.tar --format oci-archive "gotemplate:server"
	sudo ctr -n k8s.io images import /tmp/gotemplate_server.tar

k8s-deploy: clear-k8s build-k8s
	kubectl apply -f k8s/pod.yaml
	kubectl apply -f k8s/service.yaml

k8s: clear-k8s build-k8s deploy-k8s

k8s-client: export SERVER=$(shell kubectl get service server -o custom-columns=IP:.spec.clusterIP --no-headers):5051
k8s-client:
	go run cmd/client/main.go

k8s-health: export SERVER=$(shell kubectl get service server -o custom-columns=IP:.spec.clusterIP --no-headers):5051
k8s-health:
	grpc-health-probe -addr $(SERVER)
	grpc-health-probe -addr $(SERVER) -service $(SERVICE)

migrate:
	migrate -database $(DB_URI)?sslmode=disable -path $(MIGRATE_DIR) up

migrate-create:
	mkdir -p $(MIGRATE_DIR)
	migrate create -ext sql -dir $(MIGRATE_DIR) -seq $(name)

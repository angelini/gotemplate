.PHONY: install build server client health

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

server: export PORT=:5051
server: export DB_URI=postgres://postgres@10.1.1.3:5432/postgres
server:
	go run server/main.go

client: export SERVER=localhost:5051
client:
	go run client/main.go

health: export SERVER=localhost:5051
health:
	grpc-health-probe -addr $(SERVER)
	grpc-health-probe -addr $(SERVER) -service gotemplate.server.Example

.PHONY: clear-k8s build-k8s deploy-k8s client-k8s health-k8s

clear-k8s:
	kubectl delete --all service
	kubectl delete --all pod --grace-period 0 --force

build-k8s: build
	podman build -f Dockerfile -t "gotemplate:server"
	podman save -o /tmp/gotemplate_server.tar --format oci-archive "gotemplate:server"
	sudo ctr -n k8s.io images import /tmp/gotemplate_server.tar

deploy-k8s: clear-k8s build-k8s
	kubectl apply -f k8s/pod.yaml
	kubectl apply -f k8s/service.yaml

k8s: clear-k8s build-k8s deploy-k8s

client-k8s: export SERVER=$(shell kubectl get service server -o custom-columns=IP:.spec.clusterIP --no-headers):5051
client-k8s:
	go run client/main.go

health-k8s: export SERVER=$(shell kubectl get service server -o custom-columns=IP:.spec.clusterIP --no-headers):5051
health-k8s:
	grpc-health-probe -addr $(SERVER)
	grpc-health-probe -addr $(SERVER) -service gotemplate.server.Example

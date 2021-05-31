# Go GRPC Server Template

## Setup

### Requirements

- Go 1.15
- Postgresql

### Install Go Tools

```
$ make install
```

### Build

```
$ make build
$ make migrate
```

## Local

Server process

```
$ make server
```

Client process

```
$ make client
```

## K8S

Expects Containerd and Podman

### Requriements

- kubectl
- podman
- ctr

### Build

```
$ make k8s
```

### Client

```
$ make client-k8s
```

## API Testing

```
$ make test
```

## References

- https://github.com/golang-standards/project-layout

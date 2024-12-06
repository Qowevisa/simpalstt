FROM golang:1.23 as builder

WORKDIR /app/server

COPY worker/ /app/worker/

COPY worker_server/go.mod worker_server/go.sum ./

# Download grpc_health_probe binary
RUN wget -qO /usr/local/bin/grpc_health_probe https://github.com/grpc-ecosystem/grpc-health-probe/releases/download/v0.4.14/grpc_health_probe-linux-amd64 && \
    chmod +x /usr/local/bin/grpc_health_probe

RUN go mod download

COPY worker_server/. ./

RUN CGO_ENABLED=0 GOOS=linux go build -o /app/worker_server ./src

CMD ["/app/worker_server"]

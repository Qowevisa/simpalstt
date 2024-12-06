FROM golang:1.23 as builder

WORKDIR /app/graphql

COPY worker/ /app/worker/

COPY graphql_api/go.mod graphql_api/go.sum ./

RUN go mod download

COPY graphql_api/ ./

RUN CGO_ENABLED=0 GOOS=linux go build -o /app/graphql_api

CMD ["/app/graphql_api"]

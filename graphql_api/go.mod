module graphql_api

go 1.23.3

require (
	github.com/99designs/gqlgen v0.17.57
	github.com/vektah/gqlparser/v2 v2.5.19
	google.golang.org/grpc v1.68.1
	worker v0.0.0
)

require (
	github.com/agnivade/levenshtein v1.1.1 // indirect
	github.com/go-viper/mapstructure/v2 v2.2.1 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/gorilla/websocket v1.5.0 // indirect
	github.com/hashicorp/golang-lru/v2 v2.0.7 // indirect
	github.com/sosodev/duration v1.3.1 // indirect
	golang.org/x/net v0.29.0 // indirect
	golang.org/x/sys v0.25.0 // indirect
	golang.org/x/text v0.19.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20240903143218-8af14fe29dc1 // indirect
	google.golang.org/protobuf v1.35.2 // indirect
)

replace worker => ../worker/

# Simpals Test Assignment

This repository contains the implementation of the test assignment for `Simpals`. The project consists of three main components:

1. **Worker**: Located in the `worker_server` directory, responsible for processing tasks.
2. **gRPC Service**: Located in the `worker_client` directory, acts as a gRPC-based intermediary between the worker and the API.
3. **GraphQL API**: Located in the `graphql_api` directory, serves as the API layer for external access.

## Prerequisites

To run and test the project, you will need the following installed:

- [Docker](https://www.docker.com/) for container orchestration.
- [Make](https://www.gnu.org/software/make/) for managing build and generation tasks.
- [protoc](https://grpc.io/docs/protoc-installation/) for compiling protobuf files.
- [protoc-gen-go](https://pkg.go.dev/google.golang.org/protobuf/cmd/protoc-gen-go) and [protoc-gen-go-grpc](https://pkg.go.dev/google.golang.org/grpc/cmd/protoc-gen-go-grpc) for generating gRPC Go code.

## Getting Started

1. **Generate Protobuf and gRPC files**:
   Run the following command to generate necessary gRPC files:
   ```bash
   make gen_proto
   ```
2. **Start the project:**
   Use Docker Compose to start all the components. You can run it in detached mode:
   ```
   docker-compose up -d
   ```
   (or without the `-d` flag to view the logs in terminal)
3. **Check container health:**
   Each container includes health checks, so you can verify the status with:
   ```
   docker ps
   ```
   Containers that are running correctly will show (healthy) in the STATUS column.

## Running tests
To verify that all components are working correctly, you can run the provided Bash script:
```bash
./test_all_components.sh
```

This script will execute tests for all three components (`worker_server`, `worker_client`, and `graphql_api`).

## Notes

- This project has been tested on LXC containers running Ubuntu 20.04. See lxc_20.04.sh
- The system assumes you have proper Docker and Makefile support configured.
- You can find useful pre-made queries in INSTRUCTIONS.md file

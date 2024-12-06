#!/bin/bash

# Exit on any error
set -e

# List of directories to test
DIRS=("graphql_api" "worker_client" "worker_server")

# Loop through each directory and run tests
for dir in "${DIRS[@]}"; do
  if [ -d "$dir" ]; then
    echo "Running tests in $dir..."
    (cd "$dir" && go test ./... -v -cover)
  else
    echo "Directory $dir does not exist, skipping..."
  fi
done

#!/bin/bash

# Specify the path of your .proto files
PROTO_SRC="service/"

# Find all .proto files in the PROTO_SRC directory
PROTO_FILES=$(find "${PROTO_SRC}" -name "*.proto")

# Generate Go code for your .proto files
protoc --go_out=. --go_opt=paths=source_relative \
  --go-grpc_out=. --go-grpc_opt=paths=source_relative \
  $PROTO_FILES

#!/bin/bash

# Specify the path of your .proto files
PROTO_FILE="service/vectra.proto"

# Generate Go code for your .proto files
protoc --go_out=. --go_opt=paths=source_relative \
  --go-grpc_out=. --go-grpc_opt=paths=source_relative \
  $PROTO_FILE

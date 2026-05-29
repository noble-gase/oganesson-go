#!/bin/bash

set -ve

# proto
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install github.com/noble-gase/oganesson/cmd/protoc-gen-og@latest

# tag
go install github.com/favadi/protoc-go-inject-tag@latest

# swagger
go install github.com/go-swagger/go-swagger/cmd/swagger@latest
go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@latest

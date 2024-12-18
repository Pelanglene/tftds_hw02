#!/bin/bash

mkdir -p api/proto

protoc --experimental_allow_proto3_optional \
       --go_out=. \
       --go_opt=paths=source_relative \
       --go-grpc_out=. \
       --go-grpc_opt=paths=source_relative \
       api/proto/raft.proto
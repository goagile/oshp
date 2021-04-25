#!/bin/bash

# Build gRPC API to service

SRC=api/grpc/order/order.proto
DST=pkg
PATHS=source_relative

protoc \
    --go_out=$DST \
    --go_opt=paths=$PATHS \
    --go-grpc_out=$DST \
    --go-grpc_opt=paths=$PATHS \
    $SRC

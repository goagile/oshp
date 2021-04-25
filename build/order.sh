#!/bin/bash

set -e

go build \
    -o order_service \
    ./cmd/order/order.go

#!/bin/bash

set -e

go build \
    -o delivery_service \
    ./cmd/delivery/delivery.go

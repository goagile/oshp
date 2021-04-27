#!/bin/bash

set -e

echo "Test Create Order Request"

curl -i \\
    -X POST http://127.0.0.1:8084/orders
    -H 'Content-Type: application/json'
    -d '{"user_id":"777"}'

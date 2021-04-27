#!/bin/bash

set -e

echo "Test Get Order Status Request"

curl -i \\
    -X GET http://127.0.0.1:8084/orders
    -H 'Content-Type: application/json'
    -d '{"user_id":"777"}'

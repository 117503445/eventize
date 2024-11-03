#!/usr/bin/env bash

set -e
CGO_ENABLED=0 go build -o /workspace .
echo "Build successful"
#!/usr/bin/env bash

set -evx

protoc --twirp_out=. --go_out=. src/common/service.proto

# npx twirpscript
cd src/common
twirpscript service.proto
cd ../..

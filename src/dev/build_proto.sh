#!/usr/bin/env bash

protoc --twirp_out=. --go_out=. src/common/service.proto

# npx twirpscript
twirpscript src/common/service.proto

#!/usr/bin/env bash

echo "### Test all"
  set -x
  go mod tidy
  go mod vendor
  golangci-lint run ./...
  go test -v ./...

echo "### Run producer locally"
    LISTEN="localhost:8080" \
    gin -i \
      --port 28080 \
      --appPort 8080 \
      --build bin/run \
      --all \
      --excludeDir ./var \
      --excludeDir ./vendor \
      --excludeDir ./.git \
      run
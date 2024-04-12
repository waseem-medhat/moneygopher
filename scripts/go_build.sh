#!/bin/bash

for service in "$@"; do
    echo "Building binary for service: $service"
    cd services/${service}
    CGO_ENABLED=0 go build -o bin/${service} ./cmd/
    cd ../..
done

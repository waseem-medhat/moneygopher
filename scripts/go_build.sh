#!/bin/bash

cd services

for service_dir in */; do
    service=${service_dir%/}
    echo "Building binary for $service"
    cd ${service}
    CGO_ENABLED=0 go build -o bin/${service} ./cmd/
    cd ..
done

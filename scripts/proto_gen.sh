#!/bin/bash

protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    proto/money.proto \
    services/transactions/transactions.proto \
    services/accounts/accounts.proto

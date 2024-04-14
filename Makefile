ui:
	grpcui -plaintext \
		-proto proto/money.proto \
		-proto services/transactions/transactions.proto \
		-proto services/accounts/accounts.proto \
		localhost:${port}

letsgo:
	./scripts/go_build.sh
	docker compose build
	docker compose up

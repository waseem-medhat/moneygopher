ui:
	grpcui -plaintext -proto proto/money.proto -proto transactions/transactions.proto localhost:${port}

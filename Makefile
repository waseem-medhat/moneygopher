protogen:
	protoc --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		proto/money.proto \
		transactions/transactions.proto

serve:
	go run ./server

ui:
	grpcui -plaintext -proto proto/money.proto -proto transactions/transactions.proto localhost:8080

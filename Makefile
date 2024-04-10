protogen:
	protoc --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		proto/money.proto \
		transactions/transactions.proto

serve:
	go run ./${service}/cmd/

go_build:
	cd ${service} && CGO_ENABLED=0 go build -o bin/${service} ./cmd/

docker_build:
	cd ${service} && docker build -t wipdev/moneygopher-${service}:latest .

api:
	go run ./gateway/

ui:
	grpcui -plaintext -proto proto/money.proto -proto transactions/transactions.proto localhost:${port}

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

go_build_all:
	make go_build service=transactions && \
	make go_build service=gateway

docker_build_all:
	make docker_build service=transactions && \
	make docker_build service=gateway

api:
	go run ./gateway/

ui:
	grpcui -plaintext -proto proto/money.proto -proto transactions/transactions.proto localhost:${port}

package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/wipdev-tech/moneygopher/transactions"
	"google.golang.org/genproto/googleapis/type/money"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", grpcHandler)
	server := http.Server{
		Addr:    ":8080",
		Handler: mux,
	}
	fmt.Println("Listening on port 8080")
	server.ListenAndServe()
}

func grpcHandler(w http.ResponseWriter, r *http.Request) {
	opts := []grpc.DialOption{}
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))

	conn, err := grpc.Dial("0.0.0.0:8081", opts...)
	if err != nil {
		fmt.Println("failed to dial grpc:", err)
	}
	defer conn.Close()

	trxClient := transactions.NewTransactionsClient(conn)

	ctx := context.Background()
	resp, err := trxClient.Deposit(
		ctx,
		&transactions.DepositRequest{
			AccountID: "",
			Amount:    &money.Money{CurrencyCode: "USD", Units: 1000},
			Otp:       1234,
		},
	)
	if err != nil {
		fmt.Println("failed to run RPC:", err)
	} else {
		fmt.Println(resp.NewBalance)
	}
}

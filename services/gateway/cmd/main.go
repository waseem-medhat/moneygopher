package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/wipdev-tech/moneygopher/services/accounts"
	"github.com/wipdev-tech/moneygopher/services/transactions"
	"google.golang.org/genproto/googleapis/type/money"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", grpcHandler)
	mux.HandleFunc("/accounts/create", handleCreateAccount)
	server := http.Server{
		Addr:    ":8080",
		Handler: mux,
	}
	fmt.Println("Gateway is up on port 8080")
	server.ListenAndServe()
}

func grpcHandler(w http.ResponseWriter, r *http.Request) {
	opts := []grpc.DialOption{}
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))

	conn, err := grpc.Dial("transactions:8081", opts...)
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

func handleCreateAccount(w http.ResponseWriter, r *http.Request) {
	opts := []grpc.DialOption{}
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))

	conn, err := grpc.Dial("accounts:8082", opts...)
	if err != nil {
		fmt.Println("failed to dial grpc:", err)
	}
	defer conn.Close()

	accClient := accounts.NewAccountsClient(conn)

	ctx := context.Background()
	newID := uuid.NewString()
	resp, err := accClient.CreateAccount(ctx, &accounts.CreateAccountRequest{Id: newID})
	if err != nil {
		fmt.Println("failed to run RPC:", err)
	} else {
		fmt.Println("created account", resp.Id)
	}
}

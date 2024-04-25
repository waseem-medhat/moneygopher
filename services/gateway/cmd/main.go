// The API gateway is the frontend-facing service which communicates with the
// relevant services for each endpoint.
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/wipdev-tech/moneygopher/services/transactions"
	"google.golang.org/genproto/googleapis/type/money"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var insecureOpts = []grpc.DialOption{
	grpc.WithTransportCredentials(insecure.NewCredentials()),
}

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/", handleNotFound)

	mux.HandleFunc("POST /accounts", handleAccountsPost)
	mux.HandleFunc("GET /accounts/{accountID}", handleAccountsGet)

	mux.HandleFunc("POST /otps", handleOTPsPost)
	mux.HandleFunc("GET /otps", handleOTPsGet)

	server := http.Server{
		Addr:    ":" + os.Getenv("GATEWAY_PORT"),
		Handler: mux,
	}
	fmt.Println("Gateway is up on port", os.Getenv("GATEWAY_PORT"))
	server.ListenAndServe()
}

func grpcHandler(_ http.ResponseWriter, _ *http.Request) {
	conn, err := grpc.Dial("transactions:"+os.Getenv("TRANSACTIONS_PORT"), insecureOpts...)
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

func handleNotFound(w http.ResponseWriter, _ *http.Request) {
	respondError(w, http.StatusNotFound, "not found")
}

func respondJSON(w http.ResponseWriter, statusCode int, payload any) {
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(payload)
}

func respondError(w http.ResponseWriter, statusCode int, message string) {
	type err struct {
		Error string `json:"error"`
	}
	respondJSON(w, statusCode, err{Error: message})
}

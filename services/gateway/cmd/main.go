// The API gateway is the frontend-facing service which communicates with the
// relevant services for each endpoint.
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/wipdev-tech/moneygopher/services/accounts"
	"github.com/wipdev-tech/moneygopher/services/otps"
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

func handleOTPsPost(_ http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	accountID := r.URL.Query().Get("accountID")

	accoutsConn, err := grpc.Dial("accounts:"+os.Getenv("ACCOUNTS_PORT"), insecureOpts...)
	if err != nil {
		fmt.Println("failed to dial grpc:", err)
	}
	defer accoutsConn.Close()

	accountsClient := accounts.NewAccountsClient(accoutsConn)
	acc, err := accountsClient.GetAccount(ctx, &accounts.GetAccountRequest{Id: accountID})
	if err != nil {
		fmt.Println("failed to dial grpc:", err)
	} else {
		fmt.Println("found account:", acc.Id)
	}

	otpsConn, err := grpc.Dial("otp:"+os.Getenv("OTPS_PORT"), insecureOpts...)
	if err != nil {
		fmt.Println("failed to dial grpc:", err)
	}
	defer otpsConn.Close()

	otpsClient := otps.NewOtpsClient(otpsConn)

	resp, err := otpsClient.GenerateOtp(ctx, &otps.GenerateOtpRequest{AccountId: acc.Id})
	if err != nil {
		fmt.Println("failed to run RPC:", err)
	} else {
		fmt.Println("sending OTP", resp.Otp, "to number", acc.PhoneNumber)
	}
}

func handleOTPsGet(_ http.ResponseWriter, r *http.Request) {
	otp := r.URL.Query().Get("otp")

	conn, err := grpc.Dial("otp:"+os.Getenv("OTPS_PORT"), insecureOpts...)
	if err != nil {
		fmt.Println("failed to dial grpc:", err)
	}
	defer conn.Close()

	otpClient := otps.NewOtpsClient(conn)
	resp, err := otpClient.CheckOtp(r.Context(), &otps.CheckOtpRequest{Otp: otp})
	if err != nil {
		fmt.Println("failed to run RPC:", err)
	} else {
		fmt.Println("new OTP:", resp.IsValid)
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

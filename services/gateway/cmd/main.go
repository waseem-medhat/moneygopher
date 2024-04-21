package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/google/uuid"
	"github.com/wipdev-tech/moneygopher/services/accounts"
	"github.com/wipdev-tech/moneygopher/services/otp"
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
	mux.HandleFunc("/", grpcHandler)
	mux.HandleFunc("/accounts/create", handleCreateAccount)
	mux.HandleFunc("/otps/generate", handleGenerateOTP)
	server := http.Server{
		Addr:    ":" + os.Getenv("GATEWAY_PORT"),
		Handler: mux,
	}
	fmt.Println("Gateway is up on port", os.Getenv("GATEWAY_PORT"))
	server.ListenAndServe()
}

func grpcHandler(w http.ResponseWriter, r *http.Request) {
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

func handleCreateAccount(w http.ResponseWriter, r *http.Request) {
	conn, err := grpc.Dial("accounts:"+os.Getenv("ACCOUNTS_PORT"), insecureOpts...)
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

func handleGenerateOTP(w http.ResponseWriter, r *http.Request) {
	conn, err := grpc.Dial("otp:"+os.Getenv("OTP_PORT"), insecureOpts...)
	if err != nil {
		fmt.Println("failed to dial grpc:", err)
	}
	defer conn.Close()

	otpClient := otp.NewOtpClient(conn)

	ctx := context.Background()
	testID := uuid.NewString()
	resp, err := otpClient.GenerateOTP(ctx, &otp.GenerateOtpRequest{AccountId: testID})
	if err != nil {
		fmt.Println("failed to run RPC:", err)
	} else {
		fmt.Println("new OTP:", resp.Password)
	}
}

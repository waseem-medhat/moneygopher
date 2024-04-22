package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/google/uuid"
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

	mux.HandleFunc("/", grpcHandler)

	mux.HandleFunc("POST /accounts", handleAccountsPost)
	mux.HandleFunc("GET /accounts", handleAccountsGet)

	mux.HandleFunc("POST /otps", handleOTPsPost)
	mux.HandleFunc("GET /otps", handleOTPsGet)

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

func handleAccountsPost(w http.ResponseWriter, r *http.Request) {
	type In struct {
		PhoneNumber string `json:"phone_number"`
	}

	var in In
	json.NewDecoder(r.Body).Decode(&in)

	if in.PhoneNumber == "" {
		http.Error(w, "invalid phone number", http.StatusBadRequest)
		return
	}

	conn, err := grpc.Dial("accounts:"+os.Getenv("ACCOUNTS_PORT"), insecureOpts...)
	if err != nil {
		fmt.Println("failed to dial grpc:", err)
	}
	defer conn.Close()

	accClient := accounts.NewAccountsClient(conn)

	ctx := context.Background()
	newID := uuid.NewString()
	resp, err := accClient.CreateAccount(ctx, &accounts.CreateAccountRequest{
		Id:          newID,
		PhoneNumber: in.PhoneNumber,
	})

	if err != nil {
		fmt.Println("failed to run RPC:", err)
	} else {
		fmt.Println("created account", resp.Id)
	}
}

func handleAccountsGet(w http.ResponseWriter, r *http.Request) {
	accountID := r.URL.Query().Get("accountID")
	fmt.Println(accountID)

	conn, err := grpc.Dial("accounts:"+os.Getenv("ACCOUNTS_PORT"), insecureOpts...)
	if err != nil {
		fmt.Println("failed to dial grpc:", err)
	}
	defer conn.Close()

	accClient := accounts.NewAccountsClient(conn)
	resp, err := accClient.GetAccount(
		r.Context(),
		&accounts.GetAccountRequest{Id: accountID},
	)
	if err != nil {
		fmt.Println("failed to run RPC:", err)
	} else {
		fmt.Println("found account", resp.Id, "with phone number", resp.PhoneNumber)
	}
}

func handleOTPsPost(w http.ResponseWriter, r *http.Request) {
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

	otpsConn, err := grpc.Dial("otp:"+os.Getenv("OTP_PORT"), insecureOpts...)
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

func handleOTPsGet(w http.ResponseWriter, r *http.Request) {
	otp := r.URL.Query().Get("otp")

	conn, err := grpc.Dial("otp:"+os.Getenv("OTP_PORT"), insecureOpts...)
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

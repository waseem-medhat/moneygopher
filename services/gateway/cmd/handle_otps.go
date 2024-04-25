package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/wipdev-tech/moneygopher/services/accounts"
	"github.com/wipdev-tech/moneygopher/services/otps"
	"google.golang.org/grpc"
)

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

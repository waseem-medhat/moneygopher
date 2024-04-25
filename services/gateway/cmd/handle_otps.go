package main

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"os"

	"github.com/wipdev-tech/moneygopher/services/accounts"
	"github.com/wipdev-tech/moneygopher/services/otps"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

func handleOTPsPost(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	accountID := r.URL.Query().Get("accountID")

	if accountID == "" {
		respondError(w, http.StatusBadRequest, "incomplete request params")
		return
	}

	accoutsConn, err := grpc.Dial("accounts:"+os.Getenv("ACCOUNTS_PORT"), insecureOpts...)
	if err != nil {
		fmt.Println("failed to dial accounts server:", err)
		respondError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	defer accoutsConn.Close()

	accountsClient := accounts.NewAccountsClient(accoutsConn)
	acc, err := accountsClient.GetAccount(ctx, &accounts.GetAccountRequest{Id: accountID})

	if status.Convert(err).Message() == sql.ErrNoRows.Error() {
		respondError(w, http.StatusNotFound, "not found")
		return
	}

	if err != nil {
		fmt.Println("failed to get account:", err)
		respondError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	fmt.Println("found account:", acc.Id)

	otpsConn, err := grpc.Dial("otp:"+os.Getenv("OTPS_PORT"), insecureOpts...)
	if err != nil {
		fmt.Println("failed to dial otps server:", err)
		respondError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	defer otpsConn.Close()

	otpsClient := otps.NewOtpsClient(otpsConn)
	resp, err := otpsClient.GenerateOtp(ctx, &otps.GenerateOtpRequest{AccountId: acc.Id})
	if err != nil {
		fmt.Println("failed to generate otp:", err)
		respondError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	fmt.Println("sending OTP", resp.Otp, "to number", acc.PhoneNumber)
}

func handleOTPsGet(w http.ResponseWriter, r *http.Request) {
	otp := r.URL.Query().Get("otp")

	if otp == "" {
		respondError(w, http.StatusNotFound, "not found")
		return
	}

	conn, err := grpc.Dial("otp:"+os.Getenv("OTPS_PORT"), insecureOpts...)
	if err != nil {
		fmt.Println("failed to dial otps server:", err)
		respondError(w, http.StatusInternalServerError, "internal server error")
		return
	}
	defer conn.Close()

	otpClient := otps.NewOtpsClient(conn)
	resp, err := otpClient.CheckOtp(r.Context(), &otps.CheckOtpRequest{Otp: otp})
	if err != nil {
		fmt.Println("failed to check otp:", err)
		respondError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	fmt.Println("OTP valid:", resp.IsValid)
}

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/google/uuid"
	"github.com/wipdev-tech/moneygopher/services/accounts"
	"google.golang.org/grpc"
)

type account struct {
	Id          string
	PhoneNumber string
}

func handleAccountsPost(w http.ResponseWriter, r *http.Request) {
	type In struct {
		PhoneNumber string `json:"phone_number"`
	}

	var in In
	json.NewDecoder(r.Body).Decode(&in)

	if in.PhoneNumber == "" {
		respondError(w, http.StatusBadRequest, "invalid phone number")
		return
	}

	conn, err := grpc.Dial("accounts:"+os.Getenv("ACCOUNTS_PORT"), insecureOpts...)
	if err != nil {
		fmt.Println("failed to dial grpc:", err)
		respondError(w, http.StatusInternalServerError, "internal server error")
	}
	defer conn.Close()

	accClient := accounts.NewAccountsClient(conn)

	ctx := context.Background()
	newID := uuid.NewString()
	resp, err := accClient.CreateAccount(ctx, &accounts.CreateAccountRequest{
		Id:          newID,
		PhoneNumber: in.PhoneNumber,
	})

	acc := rpcAccountToAccount(resp)

	if err != nil {
		respondError(w, http.StatusInternalServerError, "internal server error")
	} else {
		respondJSON(w, http.StatusCreated, acc)
	}
}

func handleAccountsGet(_ http.ResponseWriter, r *http.Request) {
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

func rpcAccountToAccount(rpcAcc *accounts.Account) account {
	return account{
		Id:          rpcAcc.Id,
		PhoneNumber: rpcAcc.PhoneNumber,
	}
}

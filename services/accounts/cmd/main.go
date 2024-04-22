// The accounts service handles account management and storing balance
// information.
package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net"
	"os"

	_ "github.com/tursodatabase/libsql-client-go/libsql"
	pb "github.com/wipdev-tech/moneygopher/services/accounts"
	"github.com/wipdev-tech/moneygopher/services/accounts/internal/database"
	"google.golang.org/genproto/googleapis/type/money"
	"google.golang.org/grpc"
)

type accountsServer struct {
	pb.UnimplementedAccountsServer
	db *database.Queries
}

func (s *accountsServer) GetAccount(ctx context.Context, in *pb.GetAccountRequest) (*pb.Account, error) {
	dbAcc, err := s.db.GetAccountByID(ctx, in.Id)
	if err != nil {
		return nil, err
	}

	acc := &pb.Account{
		Id:          dbAcc.ID,
		PhoneNumber: dbAcc.PhoneNumber,
		Balance: &money.Money{
			CurrencyCode: "USD",
			Units:        dbAcc.BalanceDollars,
		},
	}

	return acc, err
}

func (s *accountsServer) CreateAccount(ctx context.Context, in *pb.CreateAccountRequest) (*pb.Account, error) {
	dbAcc, err := s.db.CreateAccount(ctx, database.CreateAccountParams{
		ID:          in.Id,
		PhoneNumber: in.PhoneNumber,
	})

	response := &pb.Account{
		Id:          dbAcc.ID,
		PhoneNumber: dbAcc.PhoneNumber,
		Balance: &money.Money{
			CurrencyCode: "USD",
			Units:        0,
		},
	}
	return response, err
}

func main() {
	dbURL := os.Getenv("ACCOUNTS_DB_URL")
	dbConn, err := sql.Open("libsql", dbURL)
	if err != nil {
		log.Fatal("couldn't open DB connection:", err)
	}
	defer dbConn.Close()
	db := database.New(dbConn)

	lis, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%v", os.Getenv("ACCOUNTS_PORT")))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	var opts []grpc.ServerOption

	grpcServer := grpc.NewServer(opts...)
	pb.RegisterAccountsServer(grpcServer, &accountsServer{db: db})
	fmt.Println("Accounts service is up on port", os.Getenv("ACCOUNTS_PORT"))
	grpcServer.Serve(lis)
}

package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/joho/godotenv"
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

func (s *accountsServer) GetAccount(context.Context, *pb.GetAccountRequest) (*pb.Account, error) {
	return nil, nil
}
func (s *accountsServer) CreateAccount(ctx context.Context, in *pb.CreateAccountRequest) (*pb.Account, error) {
	acc, err := s.db.CreateAccount(ctx, in.Id)
	response := &pb.Account{
		Id: acc.ID,
		Balance: &money.Money{
			CurrencyCode: "USD",
			Units:        0,
		},
	}
	return response, err
}

func main() {
	godotenv.Load()
	dbURL := os.Getenv("ACCOUNTS_DB_URL")
	dbConn, err := sql.Open("libsql", dbURL)
	if err != nil {
		log.Fatal("couldn't open DB connection:", err)
	}
	defer dbConn.Close()
	db := database.New(dbConn)

	lis, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", 8082))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	var opts []grpc.ServerOption

	grpcServer := grpc.NewServer(opts...)
	pb.RegisterAccountsServer(grpcServer, &accountsServer{db: db})
	fmt.Println("Accounts service is up on port 8082")
	grpcServer.Serve(lis)
}

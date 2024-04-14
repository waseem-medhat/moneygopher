package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"

	iso "github.com/rmg/iso4217"
	pb "github.com/wipdev-tech/moneygopher/services/transactions"
	"google.golang.org/grpc"
)

type transactionsServer struct {
	pb.UnimplementedTransactionsServer
}

func (s *transactionsServer) Deposit(ctx context.Context, in *pb.DepositRequest) (*pb.DepositResponse, error) {
	code, _ := iso.ByName(in.Amount.CurrencyCode)
	if code == 0 {
		fmt.Println("Invalid code")
	} else {
		fmt.Printf(
			"Depositing %v %v in account %v\n",
			in.Amount.Units,
			in.Amount.CurrencyCode,
			in.AccountID,
		)
	}
	res := &pb.DepositResponse{}
	return res, nil
}

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%s", os.Getenv("TRANSACTIONS_PORT")))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	var opts []grpc.ServerOption

	grpcServer := grpc.NewServer(opts...)
	pb.RegisterTransactionsServer(grpcServer, &transactionsServer{})
	fmt.Println("Transactions service is up on port", os.Getenv("TRANSACTIONS_PORT"))
	grpcServer.Serve(lis)
}

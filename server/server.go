package main

import (
	"context"
	"fmt"
	"log"
	"net"

	iso "github.com/rmg/iso4217"
	pb "github.com/wipdev-tech/microbank/transactions"
	"google.golang.org/grpc"
)

type transactionsServer struct {
	pb.UnimplementedTransactionsServer
}

func (s *transactionsServer) Deposit(ctx context.Context, in *pb.DepositRequest) (*pb.DepositResponse, error) {
	balanceCode := in.CurrentBalance.CurrencyCode
	code, _ := iso.ByName(balanceCode)
	if code == 0 {
		fmt.Println("Invalid code")
	}
	fmt.Println(in.CurrentBalance)
	res := &pb.DepositResponse{}
	return res, nil
}

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", 8080))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	var opts []grpc.ServerOption

	grpcServer := grpc.NewServer(opts...)
	pb.RegisterTransactionsServer(grpcServer, &transactionsServer{})
	grpcServer.Serve(lis)
}

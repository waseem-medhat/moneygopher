package main

import (
	"fmt"
	"log"
	"net"

	pb "github.com/wipdev-tech/moneygopher/accounts"
	"google.golang.org/grpc"
)

type accountsServer struct {
	pb.UnimplementedAccountsServer
}

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", 8081))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	var opts []grpc.ServerOption

	grpcServer := grpc.NewServer(opts...)
	pb.RegisterAccountsServer(grpcServer, &accountsServer{})
	fmt.Println("Accounts service is up on port 8082")
	grpcServer.Serve(lis)
}

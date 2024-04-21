package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"net"
	"os"

	pb "github.com/wipdev-tech/moneygopher/services/otp"
	"google.golang.org/grpc"
)

type otpServer struct {
	pb.UnimplementedOtpServer
}

func (s *otpServer) GenerateOTP(context.Context, *pb.GenerateOtpRequest) (*pb.OtpResponse, error) {
	return &pb.OtpResponse{Password: makeOTP()}, nil
}

func makeOTP() string {
	const zeroByte = byte('0')
	otpBytes := make([]byte, 6)
	for i := range otpBytes {
		otpBytes[i] = byte(rand.Intn(10)) + zeroByte
	}
	return string(otpBytes)
}

func main() {
	// dbURL := os.Getenv("OTP_DB_URL")
	// dbConn, err := sql.Open("libsql", dbURL)
	// if err != nil {
	// 	log.Fatal("couldn't open DB connection:", err)
	// }
	// defer dbConn.Close()
	// db := database.New(dbConn)

	lis, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%v", os.Getenv("OTP_PORT")))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	var opts []grpc.ServerOption

	grpcServer := grpc.NewServer(opts...)
	pb.RegisterOtpServer(grpcServer, &otpServer{})
	fmt.Println("OTP service is up on port", os.Getenv("OTP_PORT"))
	grpcServer.Serve(lis)
}

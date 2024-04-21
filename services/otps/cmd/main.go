package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"net"
	"os"
	"time"

	pb "github.com/wipdev-tech/moneygopher/services/otps"
	"google.golang.org/grpc"
)

type otpsServer struct {
	pb.UnimplementedOtpsServer
	cache *otpCache
}

func (s *otpsServer) GenerateOtp(context.Context, *pb.GenerateOtpRequest) (*pb.GenerateOtpResponse, error) {
	newOTP := makeOTP()
	s.cache.add(newOTP)
	return &pb.GenerateOtpResponse{Otp: newOTP}, nil
}

func (s *otpsServer) CheckOtp(ctx context.Context, in *pb.CheckOtpRequest) (*pb.CheckOtpResponse, error) {
	return &pb.CheckOtpResponse{IsValid: s.cache.get(in.Otp)}, nil
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
	cache := newCache()

	go func() {
		t := time.NewTicker(time.Second)
		for range t.C {
			cache.reap()
		}
	}()

	lis, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%v", os.Getenv("OTP_PORT")))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	var opts []grpc.ServerOption

	grpcServer := grpc.NewServer(opts...)
	pb.RegisterOtpsServer(grpcServer, &otpsServer{cache: cache})
	fmt.Println("OTP service is up on port", os.Getenv("OTP_PORT"))
	grpcServer.Serve(lis)
}

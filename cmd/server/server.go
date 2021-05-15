package main

import (
	"log"
	"net"

	"github.com/andrerampanelli1/fc2-grpc-go/pb/pb"
	"github.com/andrerampanelli1/fc2-grpc-go/services"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	lis, err := net.Listen("tcp", "localhost:50051")
	if err != nil {
		log.Fatalf("COULD NOT CONNECT %v", err)
	}

	grpcSv := grpc.NewServer()

	pb.RegisterUserServiceServer(grpcSv, services.NewUserService())
	reflection.Register(grpcSv)

	if err := grpcSv.Serve(lis); err != nil {
		log.Fatalf("COULD NOT SERVE: %v", err)
	}
}

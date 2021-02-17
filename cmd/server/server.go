package main

import (
	"github.com/pedrofelli/fc2-grpc/pb"
	"github.com/pedrofelli/fc2-grpc/services"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"google.golang.org/grpc"
)

func main() {
	lis, err := net.Listen("tcp", "localhost:50051")
	if err != nil {
		log.Fatal("Cound not connected: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterUserServiceServer(grpcServer, services.NewUserService())
	reflection.Register(grpcServer)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatal("Could not server: %v", err)
	}
}
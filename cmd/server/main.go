package main

import (
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/ZaneH/keep-talking/internal/actors"
	"github.com/ZaneH/keep-talking/internal/application/services"
	grpcServer "github.com/ZaneH/keep-talking/internal/infrastructure/grpc"
	pb "github.com/ZaneH/keep-talking/internal/infrastructure/grpc/proto"
)

func main() {
	// Create actor system
	actorSystem := actors.NewActorSystem()

	// Create application services
	gameService := services.NewGameService(actorSystem)

	// Create gRPC server
	gameServer := grpcServer.NewGameServiceAdapter(gameService)

	// Start gRPC server
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterGameServiceServer(s, gameServer)

	log.Printf("Server listening at %v", lis.Addr())
	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

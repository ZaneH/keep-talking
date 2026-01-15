package main

import (
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/ZaneH/defuse.party-go/internal/actors"
	appServices "github.com/ZaneH/defuse.party-go/internal/application/services"
	"github.com/ZaneH/defuse.party-go/internal/infrastructure/adapters"
	grpcServer "github.com/ZaneH/defuse.party-go/internal/infrastructure/grpc"
	pb "github.com/ZaneH/defuse.party-go/pkg/proto"
)

func main() {
	actorSystem := actors.NewActorSystem()
	actorSystemAdapter := adapters.NewActorSystemAdapter(actorSystem)
	bombService := appServices.NewBombService(actorSystemAdapter)

	gameService := appServices.NewGameService(actorSystem, bombService)
	grpcGameServiceServer := grpcServer.NewGameServiceAdapter(gameService)

	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterGameServiceServer(s, grpcGameServiceServer)

	log.Printf("Server listening at %v", lis.Addr())
	// TODO: Remove in production
	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

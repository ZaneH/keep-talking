package grpc

import (
	"context"
	"fmt"

	"github.com/ZaneH/keep-talking/internal/application/services"
	pb "github.com/ZaneH/keep-talking/internal/infrastructure/grpc/proto"
	"google.golang.org/grpc"
)

type GameServer struct {
	pb.UnimplementedGameServiceServer
	gameService *services.GameService
}

func NewGameServiceAdapter(gameService *services.GameService) *GameServer {
	return &GameServer{gameService: gameService}
}

func (s *GameServer) CreateGame(context.Context, *pb.CreateGameRequest) (*pb.CreateGameResponse, error) {
	fmt.Println("CreateGame called")
	return nil, nil
}

func (s *GameServer) SendInput(context context.Context, i *pb.PlayerInput) (*pb.PlayerInputResult, error) {
	fmt.Println("SendInput called")
	fmt.Println("Position:", i.GetModulePosition())
	switch i.GetInput().(type) {
	case *pb.PlayerInput_CutWire:
		fmt.Println("CutWire command received")
		fmt.Println("WireIndex: ", i.GetCutWire().WireIndex)
	case *pb.PlayerInput_SubmitPassword:
		fmt.Println("SubmitPassword command received")
		fmt.Println("Password: ", i.GetSubmitPassword().Password)
	}
	return &pb.PlayerInputResult{
		ModuleId: "module-id",
		Success:  true,
	}, nil
}

func (s *GameServer) StreamGameState(*pb.GameStateRequest, grpc.ServerStreamingServer[pb.GameState]) error {
	fmt.Println("StreamGameState called")
	return nil
}

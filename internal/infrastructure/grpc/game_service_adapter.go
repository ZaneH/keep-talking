package grpc

import (
	"context"
	"fmt"

	"github.com/ZaneH/keep-talking/internal/application/command"
	"github.com/ZaneH/keep-talking/internal/application/services"
	"github.com/ZaneH/keep-talking/internal/domain/valueobject"
	pb "github.com/ZaneH/keep-talking/internal/infrastructure/grpc/proto"
	"github.com/google/uuid"
)

type GameServiceAdapter struct {
	pb.UnimplementedGameServiceServer
	gameService *services.GameService
}

func NewGameServiceAdapter(gameService *services.GameService) *GameServiceAdapter {
	return &GameServiceAdapter{gameService: gameService}
}

func (s *GameServiceAdapter) CreateGame(ctx context.Context, req *pb.CreateGameRequest) (*pb.CreateGameResponse, error) {
	createGameCmd := &command.CreateGameCommand{}

	session, err := s.gameService.CreateGameSession(ctx, createGameCmd)
	if err != nil {
		return nil, fmt.Errorf("failed to create game: %v", err)
	}

	return &pb.CreateGameResponse{
		SessionId: session.GetSessionId(),
	}, nil
}

func (s *GameServiceAdapter) SendInput(ctx context.Context, i *pb.PlayerInput) (*pb.PlayerInputResult, error) {
	var cmd interface{}
	position := i.GetModulePosition()

	switch input := i.GetInput().(type) {
	case *pb.PlayerInput_CutWire:
		cmd = &command.CutWireCommand{
			SessionId:      uuid.MustParse(i.GetSessionId()),
			ModulePosition: mapProtoPositionToDomain(position),
			WireIndex:      int(input.CutWire.WireIndex),
		}
	default:
		return nil, fmt.Errorf("unknown input type: %T", input)
	}

	result, err := s.gameService.ProcessModuleInput(ctx, cmd)
	if err != nil {
		return nil, fmt.Errorf("failed to process input: %v", err)
	}

	fmt.Println("result:", result)

	return &pb.PlayerInputResult{
		ModuleId: "module-interacted-with",
		Success:  true,
	}, nil
}

func mapProtoPositionToDomain(position *pb.ModulePosition) valueobject.ModulePosition {
	var face valueobject.Face

	switch position.Face {
	case pb.BombFace_BACK:
		face = valueobject.Back
	case pb.BombFace_FRONT:
		face = valueobject.Front
	default:
		panic(fmt.Sprintf("unknown face: %v", position.Face))
	}

	return valueobject.ModulePosition{
		Face:   face,
		Row:    int(position.Row),
		Column: int(position.Col),
	}
}

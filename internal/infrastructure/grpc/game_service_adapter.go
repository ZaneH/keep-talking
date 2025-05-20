package grpc

import (
	"context"
	"fmt"
	"log"

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

	log.Printf("Created game session with ID: %s\n", session.GetSessionID())

	return &pb.CreateGameResponse{
		SessionId: session.GetSessionID().String(),
	}, nil
}

func (s *GameServiceAdapter) SendInput(ctx context.Context, i *pb.PlayerInput) (*pb.PlayerInputResult, error) {
	sessionID, err := uuid.Parse(i.GetSessionId())
	if err != nil {
		return nil, fmt.Errorf("invalid session ID: %s", i.GetSessionId())
	}
	bombID, err := uuid.Parse(i.GetBombId())
	if err != nil {
		return nil, fmt.Errorf("invalid bomb ID: %s", i.GetBombId())
	}
	moduleID, err := uuid.Parse(i.GetModuleId())
	if err != nil {
		return nil, fmt.Errorf("invalid module ID: %s", i.GetModuleId())
	}

	var cmd command.ModuleInputCommand

	switch input := i.GetInput().(type) {
	case *pb.PlayerInput_SimpleWiresInput:
		cmd = &command.SimpleWiresInputCommand{
			BaseModuleInputCommand: command.BaseModuleInputCommand{
				SessionID: sessionID,
				BombID:    bombID,
				ModuleID:  moduleID,
			},
			WireIndex: int(input.SimpleWiresInput.WireIndex),
		}
	case *pb.PlayerInput_PasswordInput:
		switch pi := input.PasswordInput.Input.(type) {
		case *pb.PasswordInput_LetterChange:
			cmd = &command.PasswordLetterChangeCommand{
				BaseModuleInputCommand: command.BaseModuleInputCommand{
					SessionID: sessionID,
					BombID:    bombID,
					ModuleID:  moduleID,
				},
				LetterIndex: int(pi.LetterChange.LetterIndex),
				Direction:   valueobject.IncrementDecrement(pi.LetterChange.Direction),
			}
		case *pb.PasswordInput_Submit:
			cmd = &command.PasswordSubmitCommand{
				BaseModuleInputCommand: command.BaseModuleInputCommand{
					SessionID: sessionID,
					BombID:    bombID,
					ModuleID:  moduleID,
				},
			}
		default:
			return nil, fmt.Errorf("unknown password input type: %T", pi)
		}
	case *pb.PlayerInput_BigButtonInput:
		cmd = &command.BigButtonInputCommand{
			BaseModuleInputCommand: command.BaseModuleInputCommand{
				SessionID: sessionID,
				BombID:    bombID,
				ModuleID:  moduleID,
			},
			PressType: valueobject.PressType(input.BigButtonInput.PressType.String()),
		}
	default:
		return nil, fmt.Errorf("unknown input type: %T", input)
	}

	res, err := s.gameService.ProcessModuleInput(ctx, cmd)
	if err != nil {
		return nil, fmt.Errorf("failed to process input: %v", err)
	}

	fmt.Printf("Processed input for session %s: %v\n", sessionID, res)

	return &pb.PlayerInputResult{
		ModuleId: i.GetModuleId(),
		Success:  res != nil && !res.(command.BaseModuleInputCommandResult).Strike,
	}, nil
}

func (s *GameServiceAdapter) GetBombs(ctx context.Context, req *pb.GetBombsRequest) (*pb.GetBombsResponse, error) {
	gameState, err := s.gameService.GetGameSession(ctx, uuid.MustParse(req.GetSessionId()))
	if err != nil {
		return nil, fmt.Errorf("failed to get game session: %v", err)
	}

	return mapGameSessionActorToProto(gameState), nil
}

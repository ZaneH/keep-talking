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
	config := buildGameConfigFromRequest(req)

	createGameCmd := &command.CreateGameCommand{
		Config: config,
	}

	session, err := s.gameService.CreateGameSession(ctx, createGameCmd)
	if err != nil {
		return nil, fmt.Errorf("failed to create game: %v", err)
	}

	log.Printf("Created game session with ID: %s\n", session.GetSessionID())

	return &pb.CreateGameResponse{
		SessionId: session.GetSessionID().String(),
	}, nil
}

// TODO: Potentially remove/replace with GameConfig/BombConfig??
func buildGameConfigFromRequest(req *pb.CreateGameRequest) valueobject.GameConfig {
	config := valueobject.NewDefaultGameConfig()

	if req.ProbabilitySimpleWires != nil {
		config.ProbabilitySimpleWires = *req.ProbabilitySimpleWires
	}

	if req.ProbabilityPassword != nil {
		config.ProbabilityPassword = *req.ProbabilityPassword
	}

	return config
}

func (s *GameServiceAdapter) SendInput(ctx context.Context, i *pb.PlayerInput) (*pb.PlayerInputResult, error) {
	sessionID := uuid.MustParse(i.GetSessionId())
	bombID := uuid.MustParse(i.GetBombId())
	moduleID := uuid.MustParse(i.GetModuleId())

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
	default:
		return nil, fmt.Errorf("unknown input type: %T", input)
	}

	res, err := s.gameService.ProcessModuleInput(ctx, cmd)
	if err != nil {
		return nil, fmt.Errorf("failed to process input: %v", err)
	}

	fmt.Printf("Processed input for session %s: %v\n", sessionID, res)

	return &pb.PlayerInputResult{
		ModuleId: "module-interacted-with",
		Success:  true,
	}, nil
}

func mapProtoPositionToDomain(position *pb.ModulePosition) valueobject.ModulePosition {
	return valueobject.ModulePosition{
		Row:    int(position.Row),
		Column: int(position.Col),
		Face:   int(position.Face),
	}
}

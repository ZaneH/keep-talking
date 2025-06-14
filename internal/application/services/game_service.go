package services

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/ZaneH/keep-talking/internal/actors"
	"github.com/ZaneH/keep-talking/internal/application/command"
	"github.com/ZaneH/keep-talking/internal/domain/services"
	"github.com/ZaneH/keep-talking/internal/domain/valueobject"
	"github.com/google/uuid"
)

type GameService struct {
	actorSystem *actors.ActorSystem
	bombService *BombService
}

func NewGameService(actorSystem *actors.ActorSystem, bombService *BombService) *GameService {
	return &GameService{actorSystem: actorSystem, bombService: bombService}
}

func (s *GameService) CreateGameSession(cmd *command.CreateGameCommand) (*actors.GameSessionActor, error) {
	config := valueobject.NewEasyGameSessionConfig(cmd.Seed)
	rng := services.NewSeededRNGFromString(config.Seed())
	session, err := s.actorSystem.CreateGameSession(rng, config)

	if err != nil {
		log.Printf("error creating game session: %v", err)
		return nil, errors.New("failed to create game session")
	}

	for _, c := range config.BombConfigs {
		_, err = s.bombService.CreateBombInSession(rng, session.GetSessionID(), c)
	}

	if err != nil {
		return nil, errors.New("failed to create bomb in session")
	}

	return session, nil
}

func (s *GameService) GetGameSession(ctx context.Context, sessionID uuid.UUID) (*actors.GameSessionActor, error) {
	sessionActor, err := s.actorSystem.GetGameSession(sessionID)
	if err != nil {
		log.Printf("error retrieving game session: %v", err)
		return nil, errors.New("game session not found")
	}

	return sessionActor, nil
}

func (s *GameService) ProcessModuleInput(ctx context.Context, cmd command.ModuleInputCommand) (interface{}, error) {
	sessionActor, err := s.actorSystem.GetGameSession(cmd.GetSessionID())
	if err != nil {
		log.Printf("error retrieving game session: %v", err)
		return nil, errors.New("game session not found")
	}

	respChan := make(chan actors.Response, 1)

	sessionActor.Send(actors.ModuleCommandMessage{
		Command:         cmd,
		ResponseChannel: respChan,
	})

	select {
	case resp := <-respChan:
		if !resp.IsSuccess() {
			return nil, resp.Error()
		}
		return resp.(actors.SuccessResponse).Data, nil

	case <-time.After(5 * time.Second):
		return nil, errors.New("timeout processing module command")

	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

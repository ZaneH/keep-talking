package services

import (
	"context"
	"errors"
	"log"

	"github.com/ZaneH/keep-talking/internal/actors"
	"github.com/ZaneH/keep-talking/internal/application/command"
	"github.com/google/uuid"
)

type GameService struct {
	actorSystem *actors.ActorSystem
}

func NewGameService(actorSystem *actors.ActorSystem) *GameService {
	return &GameService{actorSystem: actorSystem}
}

func (s *GameService) CreateGameSession(ctx context.Context, cmd *command.CreateGameCommand) (*actors.GameSessionActor, error) {
	id := uuid.New()
	session, err := s.actorSystem.GetOrCreateGameSession(id)
	if err != nil {
		log.Printf("error creating game session: %v", err)
		return nil, errors.New("failed to create game session")
	}

	return session, nil
}

func (s *GameService) ProcessModuleInput(ctx context.Context, cmd command.ModuleInputCommand) (interface{}, error) {
	sessionActor, err := s.actorSystem.GetGameSession(cmd.GetSessionId())
	if err != nil {
		log.Printf("error retrieving game session: %v", err)
		return nil, errors.New("game session not found")
	}

	result, err := sessionActor.ProcessCommand(ctx, cmd)
	if err != nil {
		log.Printf("error processing command: %v", err)
		return nil, errors.New("failed to process command")
	}

	return result, nil
}

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

func (s *GameService) CreateGameSession(cmd *command.CreateGameCommand) (*actors.GameSessionActor, valueobject.BombConfig, error) {
	var config valueobject.GameSessionConfig
	var err error

	switch cmd.ConfigType {
	case command.ConfigTypeLevel:
		config, err = valueobject.NewGameSessionConfigFromLevel(cmd.Seed, cmd.Level)
	case command.ConfigTypeMission:
		config, err = valueobject.NewGameSessionConfigFromMission(cmd.Seed, cmd.Mission)
	case command.ConfigTypeCustom:
		if cmd.CustomConfig == nil {
			return nil, valueobject.BombConfig{}, errors.New("custom config required for custom config type")
		}
		config, err = valueobject.NewGameSessionConfigFromCustom(cmd.Seed, *cmd.CustomConfig)
	default:
		// Default to easy (level 1)
		config = valueobject.NewEasyGameSessionConfig(cmd.Seed)
	}

	if err != nil {
		return nil, valueobject.BombConfig{}, err
	}

	rng := services.NewSeededRNGFromString(config.Seed())
	session, err := s.actorSystem.CreateGameSession(rng, config)

	if err != nil {
		log.Printf("error creating game session: %v", err)
		return nil, valueobject.BombConfig{}, errors.New("failed to create game session")
	}

	for _, c := range config.BombConfigs {
		_, err = s.bombService.CreateBombInSession(rng, session.GetSessionID(), c)
	}

	if err != nil {
		return nil, valueobject.BombConfig{}, errors.New("failed to create bomb in session")
	}

	// Return the first bomb config for the response
	var bombConfig valueobject.BombConfig
	if len(config.BombConfigs) > 0 {
		bombConfig = config.BombConfigs[0]
	}

	return session, bombConfig, nil
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

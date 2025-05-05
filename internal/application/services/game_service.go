package services

import (
	"context"

	"github.com/ZaneH/keep-talking/internal/actors"
	"github.com/ZaneH/keep-talking/internal/application/command"
	"github.com/ZaneH/keep-talking/internal/application/common"
)

type GameService struct {
	actorSystem *actors.ActorSystem
}

func NewGameService(actorSystem *actors.ActorSystem) *GameService {
	return &GameService{actorSystem: actorSystem}
}

// func (s *GameService) HandleCutWireCommand(ctx context.Context, cmd *command.CutWireCommand) (*common.InputResult, error) {
// 	sessionActor, err := s.actorSystem.GetOrCreateGameSession(cmd.SessionId)
// 	if err != nil {
// 		return nil, err
// 	}

// 	result, err := sessionActor.ProcessCommand(ctx, cmd)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return result.(*common.InputResult), nil
// }

func (s *GameService) HandlePlayerInput(ctx context.Context, cmd *command.PlayerInputCommand) (*common.InputResult, error) {
	sessionActor, err := s.actorSystem.GetOrCreateGameSession(cmd.SessionId)
	if err != nil {
		return nil, err
	}

	result, err := sessionActor.ProcessCommand(ctx, cmd)
	if err != nil {
		return nil, err
	}

	return result.(*common.InputResult), nil
}

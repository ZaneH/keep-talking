package services

import (
	"context"

	"github.com/ZaneH/keep-talking/internal/actors"
	"github.com/ZaneH/keep-talking/internal/domain/entities"
	"github.com/ZaneH/keep-talking/internal/domain/services"
	"github.com/ZaneH/keep-talking/internal/domain/valueobject"
	"github.com/google/uuid"
)

type BombService struct {
	actorSystem *actors.ActorSystem
	bombFactory *services.BombFactory
}

func NewBombService(actorSystem *actors.ActorSystem) *BombService {
	return &BombService{
		actorSystem: actorSystem,
		bombFactory: &services.BombFactory{},
	}
}

func (s *BombService) CreateBombInSession(
	ctx context.Context,
	sessionID uuid.UUID,
	config valueobject.BombConfig,
) (*entities.Bomb, error) {
	sessionActor, err := s.actorSystem.GetGameSession(sessionID)
	if err != nil {
		return nil, err
	}

	bomb := s.bombFactory.CreateBomb(config)

	respChan := make(chan actors.Response, 1)

	sessionActor.Send(actors.AddBombMessage{
		Bomb:            *bomb,
		ResponseChannel: respChan,
	})

	select {
	case resp := <-respChan:
		if !resp.IsSuccess() {
			return nil, resp.Error()
		}
		return bomb, nil

	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

package services

import (
	"context"

	"github.com/ZaneH/keep-talking/internal/application/ports"
	"github.com/ZaneH/keep-talking/internal/domain/entities"
	"github.com/ZaneH/keep-talking/internal/domain/valueobject"
	"github.com/google/uuid"
)

type BombService struct {
	sessionManager ports.GameSessionManager
	bombFactory    ports.BombFactory
}

func NewBombService(sessionManager ports.GameSessionManager, bombFactory ports.BombFactory) *BombService {
	return &BombService{
		sessionManager: sessionManager,
		bombFactory:    bombFactory,
	}
}

func (s *BombService) CreateBombInSession(
	ctx context.Context,
	sessionID uuid.UUID,
	config valueobject.BombConfig,
) (*entities.Bomb, error) {
	sessionActor, err := s.sessionManager.GetGameSession(sessionID)
	if err != nil {
		return nil, err
	}

	bomb := s.bombFactory.CreateBomb(config)
	if err := sessionActor.AddBomb(bomb); err != nil {
		return nil, err
	}

	return bomb, nil
}

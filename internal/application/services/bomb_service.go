package services

import (
	"log"

	"github.com/ZaneH/keep-talking/internal/application/ports"
	"github.com/ZaneH/keep-talking/internal/domain/entities"
	dPorts "github.com/ZaneH/keep-talking/internal/domain/ports"
	"github.com/ZaneH/keep-talking/internal/domain/services"
	"github.com/ZaneH/keep-talking/internal/domain/valueobject"
	"github.com/google/uuid"
)

type BombService struct {
	sessionManager ports.GameSessionManager
}

func NewBombService(sessionManager ports.GameSessionManager) *BombService {
	return &BombService{
		sessionManager: sessionManager,
	}
}

func (s *BombService) CreateBombInSession(
	rng dPorts.RandomGenerator,
	sessionID uuid.UUID,
	config valueobject.BombConfig,
) (*entities.Bomb, error) {
	sessionActor, err := s.sessionManager.GetGameSession(sessionID)
	if err != nil {
		return nil, err
	}

	mf := services.NewModuleFactory(rng)
	bf := services.NewBombFactory(mf)
	bomb := bf.CreateBomb(rng, config)

	if err := sessionActor.AddBomb(bomb); err != nil {
		log.Printf("error adding bomb to session: %v", err)
		return nil, err
	}

	return bomb, nil
}

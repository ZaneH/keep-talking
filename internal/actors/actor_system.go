package actors

import (
	"errors"
	"sync"

	"github.com/google/uuid"
)

type ActorSystem struct {
	gameSessions map[uuid.UUID]*GameSessionActor
	mu           sync.RWMutex
}

func NewActorSystem() *ActorSystem {
	return &ActorSystem{
		gameSessions: make(map[uuid.UUID]*GameSessionActor),
	}
}

func (s *ActorSystem) GetOrCreateGameSession(sessionID uuid.UUID) (*GameSessionActor, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	session, exists := s.gameSessions[sessionID]
	if !exists {
		session = NewGameSessionActor(sessionID)
		s.gameSessions[sessionID] = session
	}

	return session, nil
}

func (s *ActorSystem) GetGameSession(sessionID uuid.UUID) (*GameSessionActor, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	session, exists := s.gameSessions[sessionID]
	if !exists {
		return nil, errors.New("game session not found")
	}

	return session, nil
}

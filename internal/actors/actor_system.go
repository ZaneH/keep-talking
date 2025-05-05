package actors

import (
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

func (s *ActorSystem) GetOrCreateGameSession(sessionId uuid.UUID) (*GameSessionActor, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	session, exists := s.gameSessions[sessionId]
	if !exists {
		session = NewGameSessionActor(sessionId)
		s.gameSessions[sessionId] = session
	}

	return session, nil
}

func (s *ActorSystem) GetGameSession(sessionId uuid.UUID) (*GameSessionActor, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	session, exists := s.gameSessions[sessionId]
	if !exists {
		return nil, nil
	}

	return session, nil
}

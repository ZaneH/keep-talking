package actors

import (
	"errors"
	"sync"

	"github.com/google/uuid"
)

type ActorSystem struct {
	sessions map[uuid.UUID]*GameSessionActor
	mu       sync.RWMutex
}

func NewActorSystem() *ActorSystem {
	return &ActorSystem{
		sessions: make(map[uuid.UUID]*GameSessionActor),
	}
}

func (s *ActorSystem) CreateGameSession() (*GameSessionActor, error) {
	sessionID := uuid.New()
	sessionActor := NewGameSessionActor(sessionID)
	sessionActor.Start()

	s.mu.Lock()
	s.sessions[sessionID] = sessionActor
	s.mu.Unlock()

	return sessionActor, nil
}

func (s *ActorSystem) GetGameSession(sessionID uuid.UUID) (*GameSessionActor, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	session, exists := s.sessions[sessionID]
	if !exists {
		return nil, errors.New("game session not found")
	}

	return session, nil
}

func (s *ActorSystem) StopGameSession(sessionID uuid.UUID) error {
	s.mu.Lock()
	defer s.mu.Lock()

	session, exists := s.sessions[sessionID]
	if !exists {
		return errors.New("game session not found")
	}

	session.Stop()
	delete(s.sessions, sessionID)

	return nil
}

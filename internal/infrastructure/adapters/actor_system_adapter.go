package adapters

import (
	"errors"
	"time"

	"github.com/ZaneH/keep-talking/internal/actors"
	"github.com/ZaneH/keep-talking/internal/application/ports"
	"github.com/ZaneH/keep-talking/internal/domain/entities"
	"github.com/google/uuid"
)

type ActorSystemAdapter struct {
	actorSystem *actors.ActorSystem
}

func NewActorSystemAdapter(actorSystem *actors.ActorSystem) *ActorSystemAdapter {
	return &ActorSystemAdapter{
		actorSystem: actorSystem,
	}
}

func (a *ActorSystemAdapter) GetGameSession(sessionID uuid.UUID) (ports.SessionActor, error) {
	session, err := a.actorSystem.GetGameSession(sessionID)
	if err != nil {
		return nil, err
	}

	return &SessionActorAdapter{
		sessionActor: session,
	}, nil
}

type SessionActorAdapter struct {
	sessionActor *actors.GameSessionActor
}

func (s *SessionActorAdapter) AddBomb(bomb *entities.Bomb) error {
	respChan := make(chan actors.Response, 1)

	s.sessionActor.Send(actors.AddBombMessage{
		Bomb:            bomb,
		ResponseChannel: respChan,
	})

	select {
	case resp := <-respChan:
		if !resp.IsSuccess() {
			return resp.Error()
		}
		return nil

	case <-time.After(5 * time.Second):
		return errors.New("timeout adding bomb to session")
	}
}

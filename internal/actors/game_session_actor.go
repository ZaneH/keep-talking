package actors

import (
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/ZaneH/keep-talking/internal/domain/entities"
	"github.com/ZaneH/keep-talking/internal/domain/valueobject"
	"github.com/google/uuid"
)

type GameSessionActor struct {
	BaseActor
	session    *entities.GameSession
	bombActors map[uuid.UUID]BombActor
}

func NewGameSessionActor(sessionID uuid.UUID, config valueobject.GameConfig) *GameSessionActor {
	actor := &GameSessionActor{
		BaseActor:  NewBaseActor(100),
		bombActors: make(map[uuid.UUID]BombActor),
		session:    entities.NewGameSession(sessionID, config),
	}

	return actor
}

func (g *GameSessionActor) GetBombs() map[uuid.UUID]BombActor {
	return g.bombActors
}

func (g *GameSessionActor) Start() {
	go g.processMessages()
}

func (g *GameSessionActor) GetSessionID() uuid.UUID {
	return g.session.SessionID
}

func (g *GameSessionActor) String() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Session ID: %v\n", g.GetSessionID()))
	for _, actor := range g.bombActors {
		sb.WriteString(fmt.Sprintf("%+v", actor))
		sb.WriteString("\n")
	}
	return sb.String()
}

func (g *GameSessionActor) processMessages() {
	for {
		select {
		case msg := <-g.Mailbox():
			g.handleMessage(msg)
		case <-g.Done():
			// Stop all module actors
			for _, actor := range g.bombActors {
				actor.Stop()
			}
			return
		}
	}
}

func (g *GameSessionActor) handleMessage(msg Message) {
	switch m := msg.(type) {
	case ModuleCommandMessage:
		log.Printf("Handled")
		g.handleModuleCommand(m)
	case AddBombMessage:
		g.handleAddBombCommand(m)
	default:
		log.Printf("received unhandled message type: %T", msg)
		if m, ok := msg.(RequestMessage); ok {
			m.GetResponseChannel() <- ErrorResponse{
				Err: errors.New("unsupported message type"),
			}
		}
	}
}

func (g *GameSessionActor) handleAddBombCommand(msg AddBombMessage) {
	bomb := msg.Bomb

	bombActor := NewBombActor(bomb)
	g.bombActors[bomb.ID] = *bombActor

	msg.ResponseChannel <- &SuccessResponse{Data: bomb.ID}
}

func (g *GameSessionActor) handleModuleCommand(msg ModuleCommandMessage) {
	cmd := msg.Command
	bombID := cmd.GetBombID()
	moduleID := cmd.GetModuleID()

	moduleActor, exists := g.bombActors[bombID].moduleActors[moduleID]
	if !exists {
		msg.ResponseChannel <- ErrorResponse{
			Err: errors.New("module actor not found"),
		}
		return
	}

	// Forward the command to the module actor
	moduleActor.Send(msg)
}

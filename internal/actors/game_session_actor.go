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
	session           *entities.GameSession
	modules           map[uuid.UUID]ModuleActor
	modulesByPosition map[valueobject.ModulePosition]uuid.UUID
	config            valueobject.GameConfig
}

func NewGameSessionActor(sessionID uuid.UUID, config valueobject.GameConfig) *GameSessionActor {
	actor := &GameSessionActor{
		BaseActor:         NewBaseActor(100),
		session:           entities.NewGameSession(sessionID),
		modules:           make(map[uuid.UUID]ModuleActor),
		modulesByPosition: make(map[valueobject.ModulePosition]uuid.UUID),
		config:            config,
	}

	return actor
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
	for _, module := range g.modules {
		sb.WriteString(fmt.Sprintf("%+v", module.GetModule()))
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
			for _, module := range g.modules {
				module.Stop()
			}
			return
		}
	}
}

func (g *GameSessionActor) handleMessage(msg Message) {
	switch m := msg.(type) {
	case AddModuleMessage:
		g.handleAddModule(m)
	case GetModuleMessage:
		g.handleGetModule(m)
	case ModuleCommandMessage:
		g.handleModuleCommand(m)
	default:
		if m, ok := msg.(RequestMessage); ok {
			m.GetResponseChannel() <- ErrorResponse{
				Err: errors.New("unsupported message type"),
			}
		} else {
			log.Printf("received unhandled message type: %T", msg)
		}
	}
}

func (g *GameSessionActor) handleAddModule(msg AddModuleMessage) {
	moduleID := msg.Module.GetModuleID()
	g.modules[moduleID] = msg.Module
	g.modulesByPosition[msg.Position] = moduleID

	msg.ResponseChannel <- SuccessResponse{}
}

func (g *GameSessionActor) handleGetModule(msg GetModuleMessage) {
	moduleID, exists := g.modulesByPosition[msg.Position]
	if !exists {
		msg.GetResponseChannel() <- ErrorResponse{
			Err: errors.New("module not found"),
		}
		return
	}

	moduleActor, exists := g.modules[moduleID]
	if !exists {
		msg.GetResponseChannel() <- ErrorResponse{
			Err: errors.New("module actor not found"),
		}
		return
	}

	msg.GetResponseChannel() <- SuccessResponse{
		Data: moduleActor,
	}
}

func (g *GameSessionActor) handleModuleCommand(msg ModuleCommandMessage) {
	cmd := msg.Command
	position := cmd.GetModulePosition()

	moduleId, exists := g.modulesByPosition[position]
	if !exists {
		msg.ResponseChannel <- ErrorResponse{
			Err: errors.New("module not found"),
		}
		return
	}

	moduleActor, exists := g.modules[moduleId]
	if !exists {
		msg.ResponseChannel <- ErrorResponse{
			Err: errors.New("module actor not found"),
		}
		return
	}

	// Forward the command to the module actor
	moduleActor.Send(msg)
}

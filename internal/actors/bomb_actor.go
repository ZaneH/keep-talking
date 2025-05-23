package actors

import (
	"errors"
	"log"

	"github.com/ZaneH/keep-talking/internal/domain/entities"
	"github.com/google/uuid"
)

type BombActor struct {
	BaseActor
	bomb         *entities.Bomb
	moduleActors map[uuid.UUID]ModuleActor
}

func NewBombActor(bomb *entities.Bomb) *BombActor {
	actor := &BombActor{
		BaseActor:    NewBaseActor(100),
		bomb:         bomb,
		moduleActors: make(map[uuid.UUID]ModuleActor),
	}

	return actor
}

func (b *BombActor) Start() {
	for moduleID, module := range b.bomb.Modules {
		moduleActor, err := CreateModuleActor(b.bomb, module)
		if err != nil {
			log.Println("error creating module actor, skipped")
			continue
		}

		b.moduleActors[moduleID] = moduleActor
		moduleActor.Start()
		b.bomb.StartTimer()
	}

	go b.processMessages()
}

func (b *BombActor) GetBombID() uuid.UUID {
	return b.bomb.ID
}

func (b *BombActor) GetBomb() *entities.Bomb {
	return b.bomb
}

func (b *BombActor) GetModuleActors() map[uuid.UUID]ModuleActor {
	return b.moduleActors
}

func (b *BombActor) processMessages() {
	for {
		select {
		case msg := <-b.Mailbox():
			b.handleMessage(msg)
		case <-b.Done():
			for _, moduleActor := range b.moduleActors {
				moduleActor.Stop()
			}
			return
		}
	}
}

func (b *BombActor) handleMessage(msg Message) {
	switch m := msg.(type) {
	case ModuleCommandMessage:
		b.handleModuleCommand(m)
	default:
		if reqMsg, ok := msg.(RequestMessage); ok {
			reqMsg.GetResponseChannel() <- ErrorResponse{
				Err: errors.New("unsupported message type for bomb"),
			}
		} else {
			log.Printf("received unhandled message type: %T", msg)
		}
	}
}

func (b *BombActor) handleModuleCommand(msg ModuleCommandMessage) {
	cmd := msg.Command
	moduleID := cmd.GetModuleID()

	moduleActor, exists := b.moduleActors[moduleID]
	if !exists {
		msg.GetResponseChannel() <- ErrorResponse{
			Err: errors.New("module not found in bomb"),
		}
		return
	}

	moduleActor.Send(msg)
}

package actors

import (
	"github.com/ZaneH/defuse.party-go/internal/domain/entities"
	"github.com/google/uuid"
)

type Actor interface {
	Start()
	Stop()
	Send(message Message)
}

type ModuleActor interface {
	Actor
	GetModuleID() uuid.UUID
	GetModule() entities.Module
}

type BaseActor struct {
	mailbox chan Message
	done    chan struct{}
}

func NewBaseActor(bufferSize int) BaseActor {
	return BaseActor{
		mailbox: make(chan Message, bufferSize),
		done:    make(chan struct{}),
	}
}

func (a *BaseActor) Send(message Message) {
	select {
	case a.mailbox <- message:
	case <-a.done:
		// Message dropped
	}
}

func (a *BaseActor) Stop() {
	close(a.done)
}

func (a *BaseActor) Mailbox() <-chan Message {
	return a.mailbox
}

func (a *BaseActor) Done() <-chan struct{} {
	return a.done
}

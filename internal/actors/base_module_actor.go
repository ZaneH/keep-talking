package actors

import (
	"errors"

	"github.com/ZaneH/keep-talking/internal/domain/entities"
	"github.com/google/uuid"
)

type ModuleEntity interface {
	entities.Module
	GetModuleID() uuid.UUID
}

type BaseModuleActor struct {
	BaseActor
	module     ModuleEntity
	moduleID   uuid.UUID
	handleFunc func(msg Message)
}

func NewBaseModuleActor(module ModuleEntity, bufferSize int) BaseModuleActor {
	return BaseModuleActor{
		BaseActor:  NewBaseActor(bufferSize),
		module:     module,
		moduleID:   module.GetModuleID(),
		handleFunc: nil,
	}
}

func (a *BaseModuleActor) GetModuleID() uuid.UUID {
	return a.module.GetModuleID()
}

func (a *BaseModuleActor) GetModule() entities.Module {
	return a.module
}

func (a *BaseModuleActor) Start() {
	go a.processMessages()
}

func (a *BaseModuleActor) processMessages() {
	for {
		select {
		case msg := <-a.Mailbox():
			if a.handleFunc != nil {
				a.handleFunc(msg)
			} else {
				a.handleMessage(msg)
			}
		case <-a.Done():
			return
		}
	}
}

func (a *BaseModuleActor) handleMessage(msg Message) {
	if reqMsg, ok := msg.(RequestMessage); ok {
		reqMsg.GetResponseChannel() <- ErrorResponse{
			Err: errors.New("message not handled by this module"),
		}
	}
}

func (a *BaseModuleActor) SetMessageHandler(handler func(msg Message)) {
	a.handleFunc = handler
}

package actors

import (
	"github.com/ZaneH/defuse.party-go/internal/application/command"
	"github.com/ZaneH/defuse.party-go/internal/domain/entities"
)

type Message interface {
	MessageType() string
}

type Response interface {
	IsSuccess() bool
	Error() error
}

type RequestMessage interface {
	Message
	GetResponseChannel() chan Response
}

type ModuleCommandMessage struct {
	Command         command.ModuleInputCommand
	ResponseChannel chan Response
}

func (m ModuleCommandMessage) MessageType() string {
	return "ModuleCommand"
}

func (m ModuleCommandMessage) GetResponseChannel() chan Response {
	return m.ResponseChannel
}

type AddBombMessage struct {
	Bomb            *entities.Bomb
	ResponseChannel chan Response
}

func (m AddBombMessage) MessageType() string {
	return "AddBomb"
}

func (m AddBombMessage) GetResponseChannel() chan Response {
	return m.ResponseChannel
}

type GetBombsMessage struct {
	ResponseChannel chan Response
}

func (m GetBombsMessage) MessageType() string {
	return "GetBombs"
}

func (m GetBombsMessage) GetResponseChannel() chan Response {
	return m.ResponseChannel
}

type SuccessResponse struct {
	Data interface{}
}

func (r SuccessResponse) IsSuccess() bool {
	return true
}

func (r SuccessResponse) Error() error {
	return nil
}

type ErrorResponse struct {
	Err error
}

func (r ErrorResponse) IsSuccess() bool {
	return false
}

func (r ErrorResponse) Error() error {
	return r.Err
}

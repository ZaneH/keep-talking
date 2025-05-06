package actors

import (
	"github.com/ZaneH/keep-talking/internal/application/command"
	"github.com/ZaneH/keep-talking/internal/domain/valueobject"
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

type AddModuleMessage struct {
	Module          ModuleActor
	Position        valueobject.ModulePosition
	ResponseChannel chan Response
}

func (m AddModuleMessage) MessageType() string {
	return "AddModule"
}

type GetModuleMessage struct {
	Position        valueobject.ModulePosition
	ResponseChannel chan Response
}

func (m GetModuleMessage) MessageType() string {
	return "GetModule"
}

func (m GetModuleMessage) GetResponseChannel() chan Response {
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

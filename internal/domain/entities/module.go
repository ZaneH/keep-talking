package entities

import (
	"github.com/ZaneH/defuse.party-go/internal/domain/valueobject"
	"github.com/google/uuid"
)

type Module interface {
	GetModuleID() uuid.UUID
	GetModuleState() ModuleState
	GetPosition() valueobject.ModulePosition
	GetType() valueobject.ModuleType
	String() string
	GetBomb() *Bomb
	// AddStrike()
	SetPosition(position valueobject.ModulePosition)
	SetBomb(bomb *Bomb)
}

type ModuleState interface {
	IsSolved() bool
	MarkAsSolved()
}

type BaseModule struct {
	ModuleID uuid.UUID
	Position valueobject.ModulePosition
	bomb     *Bomb
}

func (m *BaseModule) SetBomb(bomb *Bomb) {
	m.bomb = bomb
}

func (m *BaseModule) SetPosition(position valueobject.ModulePosition) {
	m.Position = position
}

func (m *BaseModule) GetBomb() *Bomb {
	return m.bomb
}

func (m *BaseModule) GetPosition() valueobject.ModulePosition {
	return m.Position
}

func (m *BaseModule) GetModuleID() uuid.UUID {
	return m.ModuleID
}

type BaseModuleState struct {
	MarkSolved bool
}

func (b BaseModuleState) IsSolved() bool {
	return b.MarkSolved
}

func (b *BaseModuleState) MarkAsSolved() {
	b.MarkSolved = true
}

package ports

import (
	"github.com/ZaneH/keep-talking/internal/domain/entities"
	"github.com/ZaneH/keep-talking/internal/domain/valueobject"
)

type BombFactory interface {
	CreateBomb(config valueobject.BombConfig) *entities.Bomb
}

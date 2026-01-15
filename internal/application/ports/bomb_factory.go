package ports

import (
	"github.com/ZaneH/defuse.party-go/internal/domain/entities"
	"github.com/ZaneH/defuse.party-go/internal/domain/valueobject"
)

type BombFactory interface {
	CreateBomb(config valueobject.BombConfig) *entities.Bomb
}

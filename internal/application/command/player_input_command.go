package command

import (
	"github.com/ZaneH/keep-talking/internal/domain/valueobject"
	"github.com/google/uuid"
)

type PlayerInputCommand struct {
	// string player_id = 1;
	// modules.ModulePosition module_position = 2;
	// oneof input {
	//   modules.CutWireInput cut_wire = 3;
	//   modules.SubmitPasswordInput submit_password = 4;
	// }
	SessionId      uuid.UUID
	ModulePosition valueobject.ModulePosition
	Input          interface{}
}

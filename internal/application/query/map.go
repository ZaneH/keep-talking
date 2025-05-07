package query

import (
	"github.com/ZaneH/keep-talking/internal/actors"
	"github.com/ZaneH/keep-talking/internal/domain/valueobject"
	pb "github.com/ZaneH/keep-talking/internal/infrastructure/grpc/proto"
	"github.com/google/uuid"
)

func MapGameSessionActorToProto(game *actors.GameSessionActor) *pb.GetBombsResponse {
	protoGameState := pb.GetBombsResponse{}

	var bombs []*pb.Bomb
	for _, bomb := range game.GetBombActors() {
		bombs = append(bombs, &pb.Bomb{
			Id:      bomb.GetBombID().String(),
			Modules: mapModulesToProto(bomb.GetModuleActors()),
		})
	}

	protoGameState.Bombs = bombs

	return &protoGameState
}

func mapModulesToProto(modules map[uuid.UUID]actors.ModuleActor) map[string]*pb.Module {
	protoModules := make(map[string]*pb.Module)
	for _, module := range modules {
		protoModules[module.GetModuleID().String()] = &pb.Module{
			Id:     module.GetModuleID().String(),
			Type:   mapTypeToProto(module.GetModule().GetType()),
			Solved: module.GetModule().IsSolved(),
		}
	}

	return protoModules
}

func mapTypeToProto(moduleType valueobject.ModuleType) pb.Module_ModuleType {
	switch moduleType {
	case valueobject.SimpleWires:
		return pb.Module_SIMPLE_WIRES
	case valueobject.Password:
		return pb.Module_PASSWORD
	default:
		return pb.Module_UNKNOWN
	}
}

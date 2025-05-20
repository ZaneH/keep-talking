package grpc

import (
	"log"

	"github.com/ZaneH/keep-talking/internal/actors"
	"github.com/ZaneH/keep-talking/internal/domain/entities"
	"github.com/ZaneH/keep-talking/internal/domain/valueobject"
	pb "github.com/ZaneH/keep-talking/internal/infrastructure/grpc/proto"
	"github.com/google/uuid"
)

func mapGameSessionActorToProto(game *actors.GameSessionActor) *pb.GetBombsResponse {
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
		protoModule := &pb.Module{
			Id:     module.GetModuleID().String(),
			Type:   mapTypeToProto(module.GetModule().GetType()),
			Solved: module.GetModule().GetModuleState().IsSolved(),
		}

		switch module.GetModule().GetType() {
		case valueobject.SimpleWires:
			simpleWiresState, ok := module.GetModule().GetModuleState().(*entities.SimpleWiresState)
			if !ok {
				log.Printf("Expected *SimpleWiresState but got different type: %T", module.GetModule().GetModuleState())
				continue
			}

			wires := make([]*pb.Wire, 0, len(simpleWiresState.Wires))
			for _, wire := range simpleWiresState.Wires {
				wires = append(wires, &pb.Wire{
					WireColor: string(wire.WireColor),
					IsCut:     wire.IsCut,
				})
			}

			protoModule.State = &pb.Module_SimpleWires{
				SimpleWires: &pb.SimpleWiresState{
					Wires: wires,
				},
			}
		default:
			log.Fatalf("Unknown module type: %v. Couldn't provide state.", module.GetModule().GetType())
		}

		protoModules[module.GetModule().String()] = protoModule
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

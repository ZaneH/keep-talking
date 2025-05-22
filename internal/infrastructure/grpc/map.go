package grpc

import (
	"log"

	"github.com/ZaneH/keep-talking/internal/actors"
	"github.com/ZaneH/keep-talking/internal/domain/entities"
	"github.com/ZaneH/keep-talking/internal/domain/valueobject"
	pb "github.com/ZaneH/keep-talking/internal/infrastructure/grpc/proto"
	"github.com/google/uuid"
)

func mapTypeToProto(moduleType valueobject.ModuleType) pb.Module_ModuleType {
	switch moduleType {
	case valueobject.SimpleWires:
		return pb.Module_SIMPLE_WIRES
	case valueobject.Password:
		return pb.Module_PASSWORD
	case valueobject.BigButton:
		return pb.Module_BIG_BUTTON
	default:
		return pb.Module_UNKNOWN
	}
}

func mapColorToProto(color valueobject.Color) pb.Color {
	switch color {
	case valueobject.Red:
		return pb.Color_RED
	case valueobject.Blue:
		return pb.Color_BLUE
	case valueobject.White:
		return pb.Color_WHITE
	case valueobject.Yellow:
		return pb.Color_YELLOW
	case valueobject.Green:
		return pb.Color_GREEN
	case valueobject.Black:
		return pb.Color_BLACK
	case valueobject.Orange:
		return pb.Color_ORANGE
	case valueobject.Pink:
		return pb.Color_PINK
	default:
		log.Printf("Unknown color: %v. Couldn't provide state.", color)
		return pb.Color_UNKNOWN
	}
}

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
	for _, actor := range modules {
		protoModule := &pb.Module{
			Id:   actor.GetModuleID().String(),
			Type: mapTypeToProto(actor.GetModule().GetType()),
			Position: &pb.ModulePosition{
				Row:  int32(actor.GetModule().GetPosition().Row),
				Col:  int32(actor.GetModule().GetPosition().Column),
				Face: int32(actor.GetModule().GetPosition().Face),
			},
			Solved: actor.GetModule().GetModuleState().IsSolved(),
		}

		switch actor.GetModule().GetType() {
		case valueobject.SimpleWires:
			simpleWiresState, ok := actor.GetModule().GetModuleState().(*entities.SimpleWiresState)
			if !ok {
				log.Printf("Expected *SimpleWiresState but got different type: %T", actor.GetModule().GetModuleState())
				continue
			}

			wires := make([]*pb.Wire, 0, len(simpleWiresState.Wires))
			for _, wire := range simpleWiresState.Wires {
				wires = append(wires, &pb.Wire{
					WireColor: mapColorToProto(wire.WireColor),
					IsCut:     wire.IsCut,
					Index:     int32(wire.Index),
				})
			}

			protoModule.State = &pb.Module_SimpleWires{
				SimpleWires: &pb.SimpleWiresState{
					Wires: wires,
				},
			}
		case valueobject.BigButton:
			bigButtonState, ok := actor.GetModule().GetModuleState().(*entities.BigButtonState)
			if !ok {
				log.Printf("Expected *BigButtonState but got different type: %T", actor.GetModule().GetModuleState())
				continue
			}

			protoModule.State = &pb.Module_BigButton{
				BigButton: &pb.BigButtonState{
					ButtonColor: mapColorToProto(bigButtonState.ButtonColor),
					Label:       bigButtonState.Label,
				},
			}

		default:
			log.Fatalf("Unknown module type: %v. Couldn't provide state.", actor.GetModule().GetType())
		}

		protoModules[actor.GetModule().GetModuleID().String()] = protoModule
	}

	return protoModules
}

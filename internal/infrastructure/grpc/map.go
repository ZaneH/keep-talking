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
	case valueobject.Clock:
		return pb.Module_CLOCK
	case valueobject.SimonSays:
		return pb.Module_SIMON_SAYS
	default:
		log.Fatalf("Unknown module type: %v. Couldn't map type to proto.", moduleType)
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
		log.Fatalf("Unknown color: %v. Couldn't provide state.", color)
		return pb.Color_UNKNOWN
	}
}

func mapGameSessionActorToProto(game *actors.GameSessionActor) *pb.GetBombsResponse {
	protoGameState := pb.GetBombsResponse{}

	var bombs []*pb.Bomb
	for _, bombActor := range game.GetBombActors() {
		bomb := bombActor.GetBomb()

		started_at := bomb.StartedAt
		var started_at_ts int32
		if started_at != nil {
			started_at_ts = int32(bomb.StartedAt.Unix())
		}

		bombs = append(bombs, &pb.Bomb{
			Id:            bomb.ID.String(),
			SerialNumber:  bomb.SerialNumber,
			TimerDuration: int32(bomb.TimerDuration.Seconds()),
			StartedAt:     started_at_ts,
			StrikeCount:   int32(bomb.StrikeCount),
			MaxStrikes:    int32(bomb.MaxStrikes),
			Modules:       mapModulesToProto(bombActor.GetModuleActors()),
			Indicators:    mapIndicatorsToProto(bomb.Indicators),
			Batteries:     int32(bomb.Batteries),
			Ports:         mapPortsToProto(bomb.Ports),
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
		}

		moduleState := actor.GetModule().GetModuleState()
		if moduleState != nil {
			protoModule.Solved = moduleState.IsSolved()
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
					Position:  int32(wire.Position),
				})
			}

			protoModule.State = &pb.Module_SimpleWiresState{
				SimpleWiresState: &pb.SimpleWiresState{
					Wires: wires,
				},
			}
		case valueobject.BigButton:
			bigButtonState, ok := actor.GetModule().GetModuleState().(*entities.BigButtonState)
			if !ok {
				log.Printf("Expected *BigButtonState but got different type: %T", actor.GetModule().GetModuleState())
				continue
			}

			protoModule.State = &pb.Module_BigButtonState{
				BigButtonState: &pb.BigButtonState{
					ButtonColor: mapColorToProto(bigButtonState.ButtonColor),
					Label:       bigButtonState.Label,
				},
			}
		case valueobject.Clock:
		case valueobject.SimonSays:
			simonSaysState, ok := actor.GetModule().GetModuleState().(*entities.SimonSaysState)
			if !ok {
				log.Printf("Expected *SimonSaysState but got different type: %T", actor.GetModule().GetModuleState())
				continue
			}

			seq := simonSaysState.DisplaySequence
			seqProto := make([]pb.Color, len(seq))
			for i, color := range seq {
				seqProto[i] = mapColorToProto(color)
			}

			protoModule.State = &pb.Module_SimonSaysState{
				SimonSaysState: &pb.SimonSaysState{
					CurrentSequence: seqProto,
				},
			}
		case valueobject.Password:
			passwordModule, ok := actor.GetModule().(*entities.PasswordModule)
			if !ok {
				log.Printf("Expected *PasswordModule but got different type: %T", actor.GetModule().GetModuleState())
				continue
			}

			protoModule.State = &pb.Module_PasswordState{
				PasswordState: &pb.PasswordState{
					Letters: passwordModule.GetCurrentGuess(),
				},
			}
		default:
			log.Fatalf("Unknown module type: %v. Couldn't provide state.", actor.GetModule().GetType())
		}

		protoModules[actor.GetModule().GetModuleID().String()] = protoModule
	}

	return protoModules
}

func mapProtoToPressType(pressType pb.PressType) valueobject.PressType {
	switch pressType {
	case pb.PressType_TAP:
		return valueobject.PressTypeTap
	case pb.PressType_HOLD:
		return valueobject.PressTypeHold
	case pb.PressType_RELEASE:
		return valueobject.PressTypeRelease
	default:
		log.Printf("Unknown press type: %v. Falling back to Tap.", pressType)
		return valueobject.PressTypeTap
	}
}

func mapIndicatorsToProto(indicators map[string]valueobject.Indicator) map[string]*pb.Indicator {
	protoIndicators := make(map[string]*pb.Indicator, len(indicators))
	for _, indicator := range indicators {
		protoIndicators[indicator.Label] = &pb.Indicator{
			Label: indicator.Label,
			Lit:   indicator.Lit,
		}
	}
	return protoIndicators
}

func mapPortsToProto(ports []valueobject.Port) []pb.Port {
	protoPorts := make([]pb.Port, 0, len(ports))
	for _, port := range ports {
		switch port {
		case valueobject.PortDVID:
			protoPorts = append(protoPorts, pb.Port_DVID)
		case valueobject.PortRCA:
			protoPorts = append(protoPorts, pb.Port_RCA)
		case valueobject.PortPS2:
			protoPorts = append(protoPorts, pb.Port_PS2)
		case valueobject.PortRJ45:
			protoPorts = append(protoPorts, pb.Port_RJ45)
		case valueobject.PortSerial:
			protoPorts = append(protoPorts, pb.Port_SERIAL)
		}
	}
	return protoPorts
}

func mapProtoToColor(color pb.Color) valueobject.Color {
	switch color {
	case pb.Color_RED:
		return valueobject.Red
	case pb.Color_BLUE:
		return valueobject.Blue
	case pb.Color_WHITE:
		return valueobject.White
	case pb.Color_YELLOW:
		return valueobject.Yellow
	case pb.Color_GREEN:
		return valueobject.Green
	case pb.Color_BLACK:
		return valueobject.Black
	case pb.Color_ORANGE:
		return valueobject.Orange
	case pb.Color_PINK:
		return valueobject.Pink
	default:
		log.Fatalf("Unknown color: %v. Couldn't provide state.", color)
		return valueobject.Red
	}
}

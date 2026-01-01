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
	case valueobject.WiresModule:
		return pb.Module_WIRES
	case valueobject.PasswordModule:
		return pb.Module_PASSWORD
	case valueobject.BigButtonModule:
		return pb.Module_BIG_BUTTON
	case valueobject.ClockModule:
		return pb.Module_CLOCK
	case valueobject.SimonModule:
		return pb.Module_SIMON
	case valueobject.KeypadModule:
		return pb.Module_KEYPAD
	case valueobject.WhosOnFirstModule:
		return pb.Module_WHOS_ON_FIRST
	case valueobject.MemoryModule:
		return pb.Module_MEMORY
	case valueobject.MorseModule:
		return pb.Module_MORSE
	case valueobject.NeedyVentGasModule:
		return pb.Module_NEEDY_VENT_GAS
	case valueobject.NeedyKnobModule:
		return pb.Module_NEEDY_KNOB
	case valueobject.MazeModule:
		return pb.Module_MAZE
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
		case valueobject.WiresModule:
			wiresState, ok := actor.GetModule().GetModuleState().(*entities.WiresState)
			if !ok {
				log.Printf("Expected *WiresState but got different type: %T", actor.GetModule().GetModuleState())
				continue
			}

			wires := make([]*pb.Wire, 0, len(wiresState.Wires))
			for _, wire := range wiresState.Wires {
				wires = append(wires, &pb.Wire{
					WireColor: mapColorToProto(wire.WireColor),
					IsCut:     wire.IsCut,
					Position:  int32(wire.Position),
				})
			}

			protoModule.State = &pb.Module_WiresState{
				WiresState: &pb.WiresState{
					Wires: wires,
				},
			}
		case valueobject.BigButtonModule:
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
		case valueobject.ClockModule:
		case valueobject.SimonModule:
			simonState, ok := actor.GetModule().GetModuleState().(*entities.SimonState)
			if !ok {
				log.Printf("Expected *SimonState but got different type: %T", actor.GetModule().GetModuleState())
				continue
			}

			seq := simonState.DisplaySequence
			seqProto := make([]pb.Color, len(seq))
			for i, color := range seq {
				seqProto[i] = mapColorToProto(color)
			}

			protoModule.State = &pb.Module_SimonState{
				SimonState: &pb.SimonState{
					CurrentSequence: seqProto,
				},
			}
		case valueobject.PasswordModule:
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
		case valueobject.KeypadModule:
			keypadState, ok := actor.GetModule().GetModuleState().(*entities.KeypadState)
			if !ok {
				log.Printf("Expected *KeypadState but got different type: %T", actor.GetModule().GetModuleState())
				continue
			}

			displayed_symbols := make([]pb.Symbol, 0, len(keypadState.DisplayedSymbols))
			for _, symbol := range keypadState.DisplayedSymbols {
				displayed_symbols = append(displayed_symbols, mapSymbolToProto(symbol))
			}

			activated_symbols := make([]pb.Symbol, 0, len(keypadState.ActivatedSymbols))
			for symbol := range keypadState.ActivatedSymbols {
				if keypadState.ActivatedSymbols[symbol] {
					activated_symbols = append(activated_symbols, mapSymbolToProto(symbol))
				}
			}

			protoModule.State = &pb.Module_KeypadState{
				KeypadState: &pb.KeypadState{
					DisplayedSymbols: displayed_symbols,
					ActivatedSymbols: activated_symbols,
				},
			}
		case valueobject.WhosOnFirstModule:
			whosOnFirstState, ok := actor.GetModule().GetModuleState().(*entities.WhosOnFirstState)
			if !ok {
				log.Printf("Expected *WhosOnFirstState but got different type: %T", actor.GetModule().GetModuleState())
				continue
			}

			protoModule.State = &pb.Module_WhosOnFirstState{
				WhosOnFirstState: &pb.WhosOnFirstState{
					ScreenWord:  whosOnFirstState.ScreenWord,
					ButtonWords: whosOnFirstState.ButtonWords,
					Stage:       int32(whosOnFirstState.Stage),
				},
			}
		case valueobject.MemoryModule:
			memoryState, ok := actor.GetModule().GetModuleState().(*entities.MemoryState)
			if !ok {
				log.Printf("Expected *MemoryState but got different type: %T", actor.GetModule().GetModuleState())
				continue
			}

			displayedNumbers := make([]int32, len(memoryState.DisplayedNumbers))
			for i, num := range memoryState.DisplayedNumbers {
				displayedNumbers[i] = int32(num)
			}

			protoModule.State = &pb.Module_MemoryState{
				MemoryState: &pb.MemoryState{
					ScreenNumber:     int32(memoryState.ScreenNumber),
					DisplayedNumbers: displayedNumbers,
					Stage:            int32(memoryState.Stage),
				},
			}
		case valueobject.MorseModule:
			morseState, ok := actor.GetModule().GetModuleState().(*entities.MorseState)
			if !ok {
				log.Printf("Expected *MorseState but got different type: %T", actor.GetModule().GetModuleState())
				continue
			}

			protoModule.State = &pb.Module_MorseState{
				MorseState: &pb.MorseState{
					DisplayedPattern:       morseState.DisplayedPattern,
					DisplayedFrequency:     morseState.DisplayedFrequency,
					SelectedFrequencyIndex: int32(morseState.SelectedFrequencyIdx),
				},
			}
		case valueobject.NeedyVentGasModule:
			needyVentGasState, ok := actor.GetModule().GetModuleState().(*entities.NeedyVentGasState)
			if !ok {
				log.Printf("Expected *NeedyVentGasState but got different type: %T", actor.GetModule().GetModuleState())
				continue
			}

			protoModule.State = &pb.Module_NeedyVentGasState{
				NeedyVentGasState: &pb.NeedyVentGasState{
					DisplayedQuestion:  needyVentGasState.DisplayedQuestion,
					CountdownStartedAt: needyVentGasState.CountdownStartedAt,
					CountdownDuration:  int32(needyVentGasState.CountdownDuration),
				},
			}
		case valueobject.NeedyKnobModule:
			needyKnobState, ok := actor.GetModule().GetModuleState().(*entities.NeedyKnobState)
			if !ok {
				log.Printf("Expected *NeedyKnobState but got different type: %T", actor.GetModule().GetModuleState())
				continue
			}

			protoModule.State = &pb.Module_NeedyKnobState{
				NeedyKnobState: &pb.NeedyKnobState{
					DisplayedPatternFirstRow:  needyKnobState.DisplayedPattern[0],
					DisplayedPatternSecondRow: needyKnobState.DisplayedPattern[1],
					CountdownStartedAt:        needyKnobState.CountdownStartedAt,
					CountdownDuration:         int32(needyKnobState.CountdownDuration),
				},
			}
		case valueobject.MazeModule:
			mazeState, ok := actor.GetModule().GetModuleState().(*entities.MazeModuleState)
			if !ok {
				log.Printf("Expected *MazeModuleState but got different type: %T", actor.GetModule().GetModuleState())
				continue
			}

			protoModule.State = &pb.Module_MazeState{
				MazeState: &pb.MazeState{
					Maze: mapMazeToProto(mazeState.VariantToMaze()),
					PlayerPosition: &pb.Point2D{
						X: int64(mazeState.PlayerPosition.X),
						Y: int64(mazeState.PlayerPosition.Y),
					},
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

func mapSymbolToProto(symbol valueobject.Symbol) pb.Symbol {
	switch symbol {
	case valueobject.Copyright:
		return pb.Symbol_COPYRIGHT
	case valueobject.FilledStar:
		return pb.Symbol_FILLEDSTAR
	case valueobject.HollowStar:
		return pb.Symbol_HOLLOWSTAR
	case valueobject.SmileyFace:
		return pb.Symbol_SMILEYFACE
	case valueobject.DoubleK:
		return pb.Symbol_DOUBLEK
	case valueobject.Omega:
		return pb.Symbol_OMEGA
	case valueobject.SquidKnife:
		return pb.Symbol_SQUIDKNIFE
	case valueobject.Pumpkin:
		return pb.Symbol_PUMPKIN
	case valueobject.HookN:
		return pb.Symbol_HOOKN
	case valueobject.Teepee:
		log.Fatalf("Teepee symbol is not implemented in proto mapping.")
	case valueobject.Six:
		return pb.Symbol_SIX
	case valueobject.SquigglyN:
		return pb.Symbol_SQUIGGLYN
	case valueobject.At:
		return pb.Symbol_AT
	case valueobject.Ae:
		return pb.Symbol_AE
	case valueobject.MeltedThree:
		return pb.Symbol_MELTEDTHREE
	case valueobject.Euro:
		return pb.Symbol_EURO
	case valueobject.Circle:
		log.Fatalf("Circle symbol is not implemented in proto mapping.")
	case valueobject.NWithHat:
		return pb.Symbol_NWITHHAT
	case valueobject.Dragon:
		return pb.Symbol_DRAGON
	case valueobject.QuestionMark:
		return pb.Symbol_QUESTIONMARK
	case valueobject.Paragraph:
		return pb.Symbol_PARAGRAPH
	case valueobject.RightC:
		return pb.Symbol_RIGHTC
	case valueobject.LeftC:
		return pb.Symbol_LEFTC
	case valueobject.Pitchfork:
		return pb.Symbol_PITCHFORK
	case valueobject.Tripod:
		log.Fatalf("Tripod symbol is not implemented in proto mapping.")
	case valueobject.Cursive:
		return pb.Symbol_CURSIVE
	case valueobject.Tracks:
		return pb.Symbol_TRACKS
	case valueobject.Balloon:
		return pb.Symbol_BALLOON
	case valueobject.WeirdNose:
		log.Fatalf("WeirdNose symbol is not implemented in proto mapping.")
	case valueobject.UpsideDownY:
		return pb.Symbol_UPSIDEDOWNY
	case valueobject.Bt:
		return pb.Symbol_BT
	default:
		log.Fatalf("Unknown symbol: %v. Couldn't provide state.", symbol)
	}

	return pb.Symbol_COPYRIGHT
}

func mapProtoToSymbol(symbol pb.Symbol) valueobject.Symbol {
	switch symbol {
	case pb.Symbol_COPYRIGHT:
		return valueobject.Copyright
	case pb.Symbol_FILLEDSTAR:
		return valueobject.FilledStar
	case pb.Symbol_HOLLOWSTAR:
		return valueobject.HollowStar
	case pb.Symbol_SMILEYFACE:
		return valueobject.SmileyFace
	case pb.Symbol_DOUBLEK:
		return valueobject.DoubleK
	case pb.Symbol_OMEGA:
		return valueobject.Omega
	case pb.Symbol_SQUIDKNIFE:
		return valueobject.SquidKnife
	case pb.Symbol_PUMPKIN:
		return valueobject.Pumpkin
	case pb.Symbol_HOOKN:
		return valueobject.HookN
	case pb.Symbol_SIX:
		return valueobject.Six
	case pb.Symbol_SQUIGGLYN:
		return valueobject.SquigglyN
	case pb.Symbol_AT:
		return valueobject.At
	case pb.Symbol_AE:
		return valueobject.Ae
	case pb.Symbol_MELTEDTHREE:
		return valueobject.MeltedThree
	case pb.Symbol_EURO:
		return valueobject.Euro
	case pb.Symbol_NWITHHAT:
		return valueobject.NWithHat
	case pb.Symbol_DRAGON:
		return valueobject.Dragon
	case pb.Symbol_QUESTIONMARK:
		return valueobject.QuestionMark
	case pb.Symbol_PARAGRAPH:
		return valueobject.Paragraph
	case pb.Symbol_RIGHTC:
		return valueobject.RightC
	case pb.Symbol_LEFTC:
		return valueobject.LeftC
	case pb.Symbol_PITCHFORK:
		return valueobject.Pitchfork
	case pb.Symbol_CURSIVE:
		return valueobject.Cursive
	case pb.Symbol_TRACKS:
		return valueobject.Tracks
	case pb.Symbol_BALLOON:
		return valueobject.Balloon
	case pb.Symbol_UPSIDEDOWNY:
		return valueobject.UpsideDownY
	case pb.Symbol_BT:
		return valueobject.Bt
	default:
		log.Fatalf("Unknown symbol: %v. Couldn't provide state.", symbol)
	}

	return valueobject.Copyright
}

func mapMazeToProto(maze valueobject.Maze) *pb.Maze {
	rows := make([]*pb.MazeRow, 6)
	for y := range 6 {
		cells := make([]*pb.MazeCell, 6)
		for x := range 6 {
			cells[x] = &pb.MazeCell{
				Right: maze.Map[y][x].Right,
				Botom: maze.Map[y][x].Bottom,
			}
		}
		rows[y] = &pb.MazeRow{
			Cells: cells,
		}
	}

	return &pb.Maze{
		Marker_1: &pb.Point2D{X: int64(maze.Marker1.X), Y: int64(maze.Marker1.Y)},
		Marker_2: &pb.Point2D{X: int64(maze.Marker2.X), Y: int64(maze.Marker2.Y)},
		Rows:     rows,
	}
}

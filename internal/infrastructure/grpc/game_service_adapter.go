package grpc

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/ZaneH/defuse.party-go/internal/application/command"
	"github.com/ZaneH/defuse.party-go/internal/application/services"
	"github.com/ZaneH/defuse.party-go/internal/domain/valueobject"
	pb "github.com/ZaneH/defuse.party-go/pkg/proto"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GameServiceAdapter struct {
	pb.UnimplementedGameServiceServer
	gameService *services.GameService
}

func NewGameServiceAdapter(gameService *services.GameService) *GameServiceAdapter {
	return &GameServiceAdapter{gameService: gameService}
}

func (s *GameServiceAdapter) CreateGame(ctx context.Context, req *pb.CreateGameRequest) (*pb.CreateGameResponse, error) {
	cmd, err := s.protoToCreateGameCommand(req)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid config: %v", err)
	}

	session, _, err := s.gameService.CreateGameSession(cmd)
	if err != nil {
		var validationErrs valueobject.ValidationErrors
		if errors.As(err, &validationErrs) {
			return nil, status.Errorf(codes.InvalidArgument, "%v", err)
		}
		return nil, fmt.Errorf("failed to create game: %v", err)
	}

	log.Printf("Created game session with ID: %s\n", session.GetSessionID())

	return &pb.CreateGameResponse{
		SessionId: session.GetSessionID().String(),
	}, nil
}

func (s *GameServiceAdapter) protoToCreateGameCommand(req *pb.CreateGameRequest) (*command.CreateGameCommand, error) {
	cmd := &command.CreateGameCommand{}

	cfg := req.GetConfig()
	if cfg == nil {
		cmd.ConfigType = command.ConfigTypeDefault
		return cmd, nil
	}

	cmd.Seed = cfg.GetSeed()

	switch c := cfg.GetConfigType().(type) {
	case *pb.GameConfig_Level:
		cmd.ConfigType = command.ConfigTypeLevel
		cmd.Level = int(c.Level.GetLevel())

	case *pb.GameConfig_Preset:
		cmd.ConfigType = command.ConfigTypeMission
		cmd.Mission = protoMissionToDomain(c.Preset.GetMission())

	case *pb.GameConfig_Custom:
		cmd.ConfigType = command.ConfigTypeCustom
		bombConfig, err := protoCustomToBombConfig(c.Custom)
		if err != nil {
			return nil, err
		}
		cmd.CustomConfig = &bombConfig

	default:
		cmd.ConfigType = command.ConfigTypeDefault
	}

	return cmd, nil
}

func protoMissionToDomain(m pb.Mission) valueobject.Mission {
	switch m {
	case pb.Mission_THE_FIRST_BOMB:
		return valueobject.MissionTheFirstBomb
	case pb.Mission_SOMETHING_OLD_SOMETHING_NEW:
		return valueobject.MissionSomethingOldSomethingNew
	case pb.Mission_DOUBLE_YOUR_MONEY:
		return valueobject.MissionDoubleYourMoney
	case pb.Mission_ONE_STEP_UP:
		return valueobject.MissionOneStepUp
	case pb.Mission_PICK_UP_THE_PACE:
		return valueobject.MissionPickUpThePace
	case pb.Mission_A_HIDDEN_MESSAGE:
		return valueobject.MissionAHiddenMessage
	case pb.Mission_SOMETHINGS_DIFFERENT:
		return valueobject.MissionSomethingsDifferent
	case pb.Mission_ONE_GIANT_LEAP:
		return valueobject.MissionOneGiantLeap
	case pb.Mission_FAIR_GAME:
		return valueobject.MissionFairGame
	case pb.Mission_PICK_UP_THE_PACE_II:
		return valueobject.MissionPickUpThePaceII
	case pb.Mission_NO_ROOM_FOR_ERROR:
		return valueobject.MissionNoRoomForError
	case pb.Mission_EIGHT_MINUTES:
		return valueobject.MissionEightMinutes
	case pb.Mission_A_SMALL_WRINKLE:
		return valueobject.MissionASmallWrinkle
	case pb.Mission_PAY_ATTENTION:
		return valueobject.MissionPayAttention
	case pb.Mission_THE_KNOB:
		return valueobject.MissionTheKnob
	case pb.Mission_MULTI_TASKER:
		return valueobject.MissionMultiTasker
	case pb.Mission_WIRES_WIRES_EVERYWHERE:
		return valueobject.MissionWiresWiresEverywhere
	case pb.Mission_COMPUTER_HACKING:
		return valueobject.MissionComputerHacking
	case pb.Mission_WHOS_ON_FIRST_CHALLENGE:
		return valueobject.MissionWhosOnFirstChallenge
	case pb.Mission_FIENDISH:
		return valueobject.MissionFiendish
	case pb.Mission_PICK_UP_THE_PACE_III:
		return valueobject.MissionPickUpThePaceIII
	case pb.Mission_ONE_WITH_EVERYTHING:
		return valueobject.MissionOneWithEverything
	case pb.Mission_PICK_UP_THE_PACE_IV:
		return valueobject.MissionPickUpThePaceIV
	case pb.Mission_JUGGLER:
		return valueobject.MissionJuggler
	case pb.Mission_DOUBLE_TROUBLE:
		return valueobject.MissionDoubleTrouble
	case pb.Mission_I_AM_HARDCORE:
		return valueobject.MissionIAmHardcore
	case pb.Mission_BLINKENLIGHTS:
		return valueobject.MissionBlinkenlights
	case pb.Mission_APPLIED_THEORY:
		return valueobject.MissionAppliedTheory
	case pb.Mission_A_MAZE_ING:
		return valueobject.MissionAMazeIng
	case pb.Mission_SNIP_SNAP:
		return valueobject.MissionSnipSnap
	case pb.Mission_RAINBOW_TABLE:
		return valueobject.MissionRainbowTable
	case pb.Mission_BLINKENLIGHTS_II:
		return valueobject.MissionBlinkenlightsII
	default:
		return valueobject.MissionUnspecified
	}
}

func protoCustomToBombConfig(custom *pb.CustomBombConfig) (valueobject.BombConfig, error) {
	config := valueobject.BombConfig{
		Timer:             time.Duration(custom.GetTimerSeconds()) * time.Second,
		MaxStrikes:        int(custom.GetMaxStrikes()),
		NumFaces:          int(custom.GetNumFaces()),
		Rows:              int(custom.GetRows()),
		Columns:           int(custom.GetColumns()),
		MinModules:        int(custom.GetMinModules()),
		MaxModulesPerFace: int(custom.GetMaxModulesPerFace()),
		MinBatteries:      int(custom.GetMinBatteries()),
		MaxBatteries:      int(custom.GetMaxBatteries()),
		MaxIndicatorCount: int(custom.GetMaxIndicatorCount()),
		PortCount:         int(custom.GetPortCount()),
	}

	// Handle explicit module list
	if len(custom.GetModules()) > 0 {
		config.ExplicitModules = make([]valueobject.ModuleSpec, len(custom.GetModules()))
		for i, spec := range custom.GetModules() {
			moduleSpec := valueobject.ModuleSpec{
				Count: int(spec.GetCount()),
			}

			if len(spec.GetPossibleTypes()) > 0 {
				moduleSpec.PossibleTypes = make([]valueobject.ModuleType, len(spec.GetPossibleTypes()))
				for j, mt := range spec.GetPossibleTypes() {
					moduleSpec.PossibleTypes[j] = protoModuleTypeToDomain(mt)
				}
			} else {
				moduleSpec.Type = protoModuleTypeToDomain(spec.GetType())
			}

			config.ExplicitModules[i] = moduleSpec
		}
	} else {
		config.ModuleTypes = valueobject.DefaultModuleWeights()
	}

	return config, nil
}

func protoModuleTypeToDomain(mt pb.Module_ModuleType) valueobject.ModuleType {
	switch mt {
	case pb.Module_WIRES:
		return valueobject.WiresModule
	case pb.Module_PASSWORD:
		return valueobject.PasswordModule
	case pb.Module_BIG_BUTTON:
		return valueobject.BigButtonModule
	case pb.Module_SIMON:
		return valueobject.SimonModule
	case pb.Module_KEYPAD:
		return valueobject.KeypadModule
	case pb.Module_WHOS_ON_FIRST:
		return valueobject.WhosOnFirstModule
	case pb.Module_MEMORY:
		return valueobject.MemoryModule
	case pb.Module_MORSE:
		return valueobject.MorseModule
	case pb.Module_NEEDY_VENT_GAS:
		return valueobject.NeedyVentGasModule
	case pb.Module_NEEDY_KNOB:
		return valueobject.NeedyKnobModule
	case pb.Module_MAZE:
		return valueobject.MazeModule
	default:
		return valueobject.WiresModule // fallback
	}
}

func domainModuleTypeToProto(mt valueobject.ModuleType) pb.Module_ModuleType {
	switch mt {
	case valueobject.ClockModule:
		return pb.Module_CLOCK
	case valueobject.WiresModule:
		return pb.Module_WIRES
	case valueobject.PasswordModule:
		return pb.Module_PASSWORD
	case valueobject.BigButtonModule:
		return pb.Module_BIG_BUTTON
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
		return pb.Module_UNKNOWN
	}
}

func (s *GameServiceAdapter) SendInput(ctx context.Context, i *pb.PlayerInput) (*pb.PlayerInputResult, error) {
	sessionID, err := uuid.Parse(i.GetSessionId())
	if err != nil {
		return nil, fmt.Errorf("invalid session ID: %s", i.GetSessionId())
	}
	bombID, err := uuid.Parse(i.GetBombId())
	if err != nil {
		return nil, fmt.Errorf("invalid bomb ID: %s", i.GetBombId())
	}
	moduleID, err := uuid.Parse(i.GetModuleId())
	if err != nil {
		return nil, fmt.Errorf("invalid module ID: %s", i.GetModuleId())
	}

	var cmd command.ModuleInputCommand

	switch input := i.GetInput().(type) {
	case *pb.PlayerInput_WiresInput:
		cmd = &command.WiresInputCommand{
			BaseModuleInputCommand: command.BaseModuleInputCommand{
				SessionID: sessionID,
				BombID:    bombID,
				ModuleID:  moduleID,
			},
			WirePosition: int(input.WiresInput.WirePosition),
		}
	case *pb.PlayerInput_PasswordInput:
		switch pi := input.PasswordInput.Input.(type) {
		case *pb.PasswordInput_LetterChange:
			cmd = &command.PasswordLetterChangeCommand{
				BaseModuleInputCommand: command.BaseModuleInputCommand{
					SessionID: sessionID,
					BombID:    bombID,
					ModuleID:  moduleID,
				},
				LetterIndex: int(pi.LetterChange.LetterIndex),
				Direction:   valueobject.IncrementDecrement(pi.LetterChange.Direction),
			}
		case *pb.PasswordInput_Submit:
			cmd = &command.PasswordSubmitCommand{
				BaseModuleInputCommand: command.BaseModuleInputCommand{
					SessionID: sessionID,
					BombID:    bombID,
					ModuleID:  moduleID,
				},
			}
		default:
			return nil, fmt.Errorf("unknown password input type: %T", pi)
		}
	case *pb.PlayerInput_BigButtonInput:
		cmd = &command.BigButtonInputCommand{
			BaseModuleInputCommand: command.BaseModuleInputCommand{
				SessionID: sessionID,
				BombID:    bombID,
				ModuleID:  moduleID,
			},
			PressType:        mapProtoToPressType(input.BigButtonInput.PressType),
			ReleaseTimestamp: input.BigButtonInput.ReleaseTimestamp,
		}
	case *pb.PlayerInput_SimonInput:
		cmd = &command.SimonInputCommand{
			BaseModuleInputCommand: command.BaseModuleInputCommand{
				SessionID: sessionID,
				BombID:    bombID,
				ModuleID:  moduleID,
			},
			Color: mapProtoToColor(input.SimonInput.Color),
		}
	case *pb.PlayerInput_KeypadInput:
		cmd = &command.KeypadInputCommand{
			BaseModuleInputCommand: command.BaseModuleInputCommand{
				SessionID: sessionID,
				BombID:    bombID,
				ModuleID:  moduleID,
			},
			Symbol: mapProtoToSymbol(input.KeypadInput.Symbol),
		}
	case *pb.PlayerInput_WhosOnFirstInput:
		cmd = &command.WhosOnFirstInputCommand{
			BaseModuleInputCommand: command.BaseModuleInputCommand{
				SessionID: sessionID,
				BombID:    bombID,
				ModuleID:  moduleID,
			},
			Word: input.WhosOnFirstInput.Word,
		}
	case *pb.PlayerInput_MemoryInput:
		cmd = &command.MemoryInputCommand{
			BaseModuleInputCommand: command.BaseModuleInputCommand{
				SessionID: sessionID,
				BombID:    bombID,
				ModuleID:  moduleID,
			},
			ButtonIndex: int(input.MemoryInput.ButtonIndex),
		}
	case *pb.PlayerInput_MorseInput:
		switch mi := input.MorseInput.Input.(type) {
		case *pb.MorseInput_FrequencyChange:
			cmd = &command.MorseChangeFrequencyCommand{
				BaseModuleInputCommand: command.BaseModuleInputCommand{
					SessionID: sessionID,
					BombID:    bombID,
					ModuleID:  moduleID,
				},
				Direction: valueobject.IncrementDecrement(mi.FrequencyChange.Direction),
			}
		case *pb.MorseInput_Tx:
			cmd = &command.MorseTxCommand{
				BaseModuleInputCommand: command.BaseModuleInputCommand{
					SessionID: sessionID,
					BombID:    bombID,
					ModuleID:  moduleID,
				},
			}
		}
	case *pb.PlayerInput_NeedyVentGasInput:
		cmd = &command.NeedyVentGasCommand{
			BaseModuleInputCommand: command.BaseModuleInputCommand{
				SessionID: sessionID,
				BombID:    bombID,
				ModuleID:  moduleID,
			},
			Input: input.NeedyVentGasInput.Input,
		}
	case *pb.PlayerInput_NeedyKnobInput:
		cmd = &command.NeedyKnobCommand{
			BaseModuleInputCommand: command.BaseModuleInputCommand{
				SessionID: sessionID,
				BombID:    bombID,
				ModuleID:  moduleID,
			},
		}
	case *pb.PlayerInput_MazeInput:
		cmd = &command.MazeCommand{
			BaseModuleInputCommand: command.BaseModuleInputCommand{
				SessionID: sessionID,
				BombID:    bombID,
				ModuleID:  moduleID,
			},
			Direction: valueobject.CardinalDirection(input.MazeInput.Direction),
		}
	default:
		return nil, fmt.Errorf("unknown input type: %T", input)
	}

	res, err := s.gameService.ProcessModuleInput(ctx, cmd)
	if err != nil {
		return nil, fmt.Errorf("failed to process input: %v", err)
	}

	fmt.Printf("Processed input for session %s: %v\n", sessionID, res)

	// Get the current bomb state to include in the response
	session, err := s.gameService.GetGameSession(ctx, sessionID)
	if err != nil {
		return nil, fmt.Errorf("failed to get game session after input: %v", err)
	}

	bombActors := session.GetBombActors()
	bombActor, exists := bombActors[bombID]
	if !exists {
		return nil, fmt.Errorf("bomb not found in session")
	}

	bomb := bombActor.GetBomb()
	bombStatus := &pb.BombStatus{
		StrikeCount: int32(bomb.StrikeCount),
		MaxStrikes:  int32(bomb.MaxStrikes),
		Exploded:    bomb.StrikeCount >= bomb.MaxStrikes,
	}

	switch cmdResult := res.(type) {
	case *command.WiresInputCommandResult:
		return &pb.PlayerInputResult{
			ModuleId:   i.GetModuleId(),
			Strike:     res != nil && cmdResult.Strike,
			Solved:     res != nil && cmdResult.Solved,
			BombStatus: bombStatus,
		}, nil
	case *command.BigButtonInputCommandResult:
		var color pb.Color
		if cmdResult.StripColor != nil {
			color = mapColorToProto(*cmdResult.StripColor)
		}

		return &pb.PlayerInputResult{
			ModuleId:   i.GetModuleId(),
			Strike:     res != nil && cmdResult.Strike,
			Solved:     res != nil && cmdResult.Solved,
			BombStatus: bombStatus,
			Result: &pb.PlayerInputResult_BigButtonInputResult{
				BigButtonInputResult: &pb.BigButtonInputResult{
					StripColor: color,
				},
			},
		}, nil
	case *command.SimonInputCommandResult:
		sequence := make([]pb.Color, len(cmdResult.DisplaySequence))
		for i, color := range cmdResult.DisplaySequence {
			sequence[i] = mapColorToProto(color)
		}

		return &pb.PlayerInputResult{
			ModuleId:   i.GetModuleId(),
			Strike:     res != nil && cmdResult.Strike,
			Solved:     res != nil && cmdResult.Solved,
			BombStatus: bombStatus,
			Result: &pb.PlayerInputResult_SimonInputResult{
				SimonInputResult: &pb.SimonInputResult{
					DisplaySequence: sequence,
					HasFinishedSeq:  cmdResult.HasFinishedSeq,
				},
			},
		}, nil
	case *command.PasswordCommandResult:
		return &pb.PlayerInputResult{
			ModuleId:   i.GetModuleId(),
			BombStatus: bombStatus,
			Strike:     res != nil && cmdResult.Strike,
			Solved:     res != nil && cmdResult.Solved,
			Result: &pb.PlayerInputResult_PasswordInputResult{
				PasswordInputResult: &pb.PasswordInputResult{
					PasswordState: &pb.PasswordState{
						Letters: cmdResult.Letters,
					},
				},
			},
		}, nil
	case *command.KeypadInputCommandResult:
		activatedSymbols := make([]pb.Symbol, 0, len(cmdResult.ActivatedSymbols))
		for sym, active := range cmdResult.ActivatedSymbols {
			if active {
				activatedSymbols = append(activatedSymbols, mapSymbolToProto(sym))
			}
		}

		displayedSymbols := make([]pb.Symbol, len(cmdResult.DisplayedSymbols))
		for i, sym := range cmdResult.DisplayedSymbols {
			displayedSymbols[i] = mapSymbolToProto(sym)
		}

		return &pb.PlayerInputResult{
			ModuleId:   i.GetModuleId(),
			BombStatus: bombStatus,
			Strike:     res != nil && cmdResult.Strike,
			Solved:     res != nil && cmdResult.Solved,
			Result: &pb.PlayerInputResult_KeypadInputResult{
				KeypadInputResult: &pb.KeypadInputResult{
					KeypadState: &pb.KeypadState{
						ActivatedSymbols: activatedSymbols,
						DisplayedSymbols: displayedSymbols,
					},
				},
			},
		}, nil
	case *command.WhosOnFirstInputCommandResult:
		return &pb.PlayerInputResult{
			ModuleId:   i.GetModuleId(),
			BombStatus: bombStatus,
			Strike:     res != nil && cmdResult.Strike,
			Solved:     res != nil && cmdResult.Solved,
			Result: &pb.PlayerInputResult_WhosOnFirstInputResult{
				WhosOnFirstInputResult: &pb.WhosOnFirstInputResult{
					WhosOnFirstState: &pb.WhosOnFirstState{
						ScreenWord:  cmdResult.ScreenWord,
						ButtonWords: cmdResult.ButtonWords,
						Stage:       int32(cmdResult.Stage),
					},
				},
			},
		}, nil
	case *command.MemoryInputCommandResult:
		int32Slice := make([]int32, len(cmdResult.DisplayedNumbers))
		for i, num := range cmdResult.DisplayedNumbers {
			int32Slice[i] = int32(num)
		}

		return &pb.PlayerInputResult{
			ModuleId:   i.GetModuleId(),
			BombStatus: bombStatus,
			Strike:     res != nil && cmdResult.Strike,
			Solved:     res != nil && cmdResult.Solved,
			Result: &pb.PlayerInputResult_MemoryInputResult{
				MemoryInputResult: &pb.MemoryInputResult{
					MemoryState: &pb.MemoryState{
						ScreenNumber:     int32(cmdResult.ScreenNumber),
						DisplayedNumbers: int32Slice,
						Stage:            int32(cmdResult.Stage),
					},
				},
			},
		}, nil
	case *command.MorseCommandResult:
		return &pb.PlayerInputResult{
			ModuleId:   i.GetModuleId(),
			BombStatus: bombStatus,
			Strike:     res != nil && cmdResult.Strike,
			Solved:     res != nil && cmdResult.Solved,
			Result: &pb.PlayerInputResult_MorseInputResult{
				MorseInputResult: &pb.MorseInputResult{
					MorseState: &pb.MorseState{
						DisplayedPattern:       cmdResult.DisplayedPattern,
						DisplayedFrequency:     cmdResult.DisplayedFrequency,
						SelectedFrequencyIndex: int32(cmdResult.SelectedFrequencyIdx),
					},
				},
			},
		}, nil
	case *command.NeedyVentGasCommandResult:
		return &pb.PlayerInputResult{
			ModuleId:   i.GetModuleId(),
			BombStatus: bombStatus,
			Strike:     res != nil && cmdResult.Strike,
			Solved:     res != nil && cmdResult.Solved,
			Result: &pb.PlayerInputResult_NeedyVentGasInputResult{
				NeedyVentGasInputResult: &pb.NeedyVentGasInputResult{
					NeedyVentGasState: &pb.NeedyVentGasState{
						DisplayedQuestion:  cmdResult.DisplayedQuestion,
						CountdownStartedAt: cmdResult.CountdownStartedAt,
						CountdownDuration:  int32(cmdResult.CountdownDuration),
					},
				},
			},
		}, nil
	case *command.NeedyKnobCommandResult:
		return &pb.PlayerInputResult{
			ModuleId:   i.GetModuleId(),
			BombStatus: bombStatus,
			Strike:     res != nil && cmdResult.Strike,
			Solved:     res != nil && cmdResult.Solved,
			Result: &pb.PlayerInputResult_NeedyKnobInputResult{
				NeedyKnobInputResult: &pb.NeedyKnobInputResult{
					NeedyKnobState: &pb.NeedyKnobState{
						DisplayedPatternFirstRow:  cmdResult.DisplayedPattern[0],
						DisplayedPatternSecondRow: cmdResult.DisplayedPattern[1],
						DialDirection:             pb.CardinalDirection(cmdResult.DialDirection),
						CountdownStartedAt:        cmdResult.CoundownStartedAt,
						CountdownDuration:         int32(cmdResult.CountdownDuration),
					},
				},
			},
		}, nil
	case *command.MazeInputCommandResult:
		return &pb.PlayerInputResult{
			ModuleId:   i.GetModuleId(),
			BombStatus: bombStatus,
			Strike:     res != nil && cmdResult.Strike,
			Solved:     res != nil && cmdResult.Solved,
			Result: &pb.PlayerInputResult_MazeInputResult{
				MazeInputResult: &pb.MazeInputResult{
					MazeState: &pb.MazeState{
						// Marker_1: *pb.Point2D,
						// Marker_2: *pb.Point2D,
						PlayerPosition: mapPoint2DToProto(cmdResult.PlayerPosition),
						GoalPosition:   mapPoint2DToProto(cmdResult.GoalPosition),
					},
				},
			},
		}, nil
	case nil:
		return nil, nil
	default:
		return nil, fmt.Errorf("unknown result type: %T", res)
	}
}

func (s *GameServiceAdapter) GetBombs(ctx context.Context, req *pb.GetBombsRequest) (*pb.GetBombsResponse, error) {
	gameState, err := s.gameService.GetGameSession(ctx, uuid.MustParse(req.GetSessionId()))
	if err != nil {
		return nil, fmt.Errorf("failed to get game session: %v", err)
	}

	return mapGameSessionActorToProto(gameState), nil
}

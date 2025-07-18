package grpc

import (
	"context"
	"fmt"
	"log"

	"github.com/ZaneH/keep-talking/internal/application/command"
	"github.com/ZaneH/keep-talking/internal/application/services"
	"github.com/ZaneH/keep-talking/internal/domain/valueobject"
	pb "github.com/ZaneH/keep-talking/internal/infrastructure/grpc/proto"
	"github.com/google/uuid"
)

type GameServiceAdapter struct {
	pb.UnimplementedGameServiceServer
	gameService *services.GameService
}

func NewGameServiceAdapter(gameService *services.GameService) *GameServiceAdapter {
	return &GameServiceAdapter{gameService: gameService}
}

func (s *GameServiceAdapter) CreateGame(ctx context.Context, req *pb.CreateGameRequest) (*pb.CreateGameResponse, error) {
	createGameCmd := &command.CreateGameCommand{}

	session, err := s.gameService.CreateGameSession(createGameCmd)
	if err != nil {
		return nil, fmt.Errorf("failed to create game: %v", err)
	}

	log.Printf("Created game session with ID: %s\n", session.GetSessionID())

	return &pb.CreateGameResponse{
		SessionId: session.GetSessionID().String(),
	}, nil
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
	default:
		return nil, fmt.Errorf("unknown input type: %T", input)
	}

	res, err := s.gameService.ProcessModuleInput(ctx, cmd)
	if err != nil {
		return nil, fmt.Errorf("failed to process input: %v", err)
	}

	fmt.Printf("Processed input for session %s: %v\n", sessionID, res)

	switch cmdResult := res.(type) {
	case *command.WiresInputCommandResult:
		return &pb.PlayerInputResult{
			ModuleId: i.GetModuleId(),
			Strike:   res != nil && cmdResult.Strike,
			Solved:   res != nil && cmdResult.Solved,
		}, nil
	case *command.BigButtonInputCommandResult:
		var color pb.Color
		if cmdResult.StripColor != nil {
			color = mapColorToProto(*cmdResult.StripColor)
		}

		return &pb.PlayerInputResult{
			ModuleId: i.GetModuleId(),
			Strike:   res != nil && cmdResult.Strike,
			Solved:   res != nil && cmdResult.Solved,
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
			ModuleId: i.GetModuleId(),
			Strike:   res != nil && cmdResult.Strike,
			Solved:   res != nil && cmdResult.Solved,
			Result: &pb.PlayerInputResult_SimonInputResult{
				SimonInputResult: &pb.SimonInputResult{
					DisplaySequence: sequence,
					HasFinishedSeq:  cmdResult.HasFinishedSeq,
				},
			},
		}, nil
	case *command.PasswordCommandResult:
		return &pb.PlayerInputResult{
			ModuleId: i.GetModuleId(),
			Strike:   res != nil && cmdResult.Strike,
			Solved:   res != nil && cmdResult.Solved,
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
			ModuleId: i.GetModuleId(),
			Strike:   res != nil && cmdResult.Strike,
			Solved:   res != nil && cmdResult.Solved,
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
			ModuleId: i.GetModuleId(),
			Strike:   res != nil && cmdResult.Strike,
			Solved:   res != nil && cmdResult.Solved,
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
			ModuleId: i.GetModuleId(),
			Strike:   res != nil && cmdResult.Strike,
			Solved:   res != nil && cmdResult.Solved,
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
			ModuleId: i.GetModuleId(),
			Strike:   res != nil && cmdResult.Strike,
			Solved:   res != nil && cmdResult.Solved,
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
			ModuleId: i.GetModuleId(),
			Strike:   res != nil && cmdResult.Strike,
			Solved:   res != nil && cmdResult.Solved,
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

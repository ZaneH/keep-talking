package interfaces

import (
	"context"

	"github.com/ZaneH/keep-talking/internal/application/command"
	"github.com/ZaneH/keep-talking/internal/application/common"
	domain "github.com/ZaneH/keep-talking/internal/domain/entities"
	"github.com/ZaneH/keep-talking/internal/infrastructure/grpc/proto"
	"google.golang.org/grpc"
)

type GameServer interface {
	HandleCommand(ctx context.Context, command interface{}) (interface{}, error)
	CreateGame(ctx context.Context, req *command.CreateGameCommand) (*domain.GameSession, error)
	SendInput(ctx context.Context, input *command.PlayerInputCommand) (*common.InputResult, error)
	StreamGameState(req *proto.GameSessionRequest, stream grpc.ServerStreamingServer[proto.GameSession]) error
}

package services

import (
	gpb "github.com/ZaneH/keep-talking/internal/infrastructure/grpc/proto"
)

type GameService struct {
	// Add any necessary fields here, such as a logger or a database connection
}

func NewGameService() *GameService {
	g := &gpb.GameServiceClient{}

	return &GameService{}
}

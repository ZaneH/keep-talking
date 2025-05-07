# Keep Talking

Keep Talking is an open-source re-implementation of one of my favorite games: "Keep Talking and Nobody Explodes".

The game is designed for two or more players, where one player (the "Defuser") is tasked with defusing a bomb while the other players (the "Experts") provide instructions from a bomb defusal manual. The catch is that the Defuser cannot see the manual, and the Experts cannot see the bomb.

This repo is a practical exercise in DDD (Domain-Driven Design) and the actor model. The game is implemented in Go and uses gRPC for interaction.

## Setup

To run the game, you need Go installed on your machine. You can download it from the [official Go website](https://go.dev/dl/).

### Clone the Repository

```bash
$ git clone https://github.com/ZaneH/keep-talking.git
$ cd keep-talking
```

### Install Dependencies

```bash
$ go mod tidy
```

### Run the Game
```bash
$ go run cmd/server/main.go
```
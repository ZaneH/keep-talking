# Keep Talking

Keep Talking is an open-source gRPC implementation of the popular game "Keep Talking and Nobody Explodes". The game is designed for two or more players, where one player (the "Defuser") is tasked with defusing a bomb while the other players (the "Experts") provide instructions from a bomb defusal manual. The catch is that the Defuser cannot see the manual, and the Experts cannot see the bomb.

## Features

- **gRPC Communication**: The game uses gRPC for communication between the Defuser and Experts, allowing for real-time interaction.
- **DDD Architecture**: The code is organized using Domain-Driven Design principles, making it easy to understand and extend.
- **Well Tested**: The code is thoroughly tested, ensuring that the game runs smoothly and without bugs.

- **Easy to Extend**: The game is designed to be easily extendable, allowing for new features and game modes to be added in the future.

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
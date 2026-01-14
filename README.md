# Keep Talking

![CI status](https://github.com/zaneh/keep-talking/actions/workflows/ci.yml/badge.svg)

Keep Talking is an open-source re-implementation of one of my favorite games: "Keep Talking and Nobody Explodes".

The game is designed for two or more players, where one player (the "Defuser") is tasked with defusing a bomb while the other players (the "Experts") provide instructions from a bomb defusal manual.
The catch is that the Defuser cannot see the manual, and the Experts cannot see the bomb.

This repo serves as a practical exercise in DDD (Domain-Driven Design) and the actor model.

## Architecture

The project follows a hexagonal architecture with clear separation of concerns:

- **Domain Layer**: Core game logic and entities
- **Application Layer**: Use cases and application services
- **Infrastructure Layer**: External adapters and implementations

The bomb defusal logic is implemented using the actor model, where different components of the game (bombs, modules, etc.) are represented as actors that communicate via messages. Multiple game sessions can run concurrently, each with its own state and actors.

## API Design

The server exposes a gRPC API and an optional HTTP Proxy using [grpc-gateway](https://github.com/grpc-ecosystem/grpc-gateway), making it easy for clients to interact with the game. Protocol Buffers (protobuf) are used for efficient, language-agnostic data serialization.
This design allows developers to build their own game clients in any language that supports gRPC.

While this implementation focuses on gRPC/HTTP, the Domain-Driven Design approach means that alternative interfaces (like WebSockets or) could be implemented without modifying the core game logic.

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

```bash
$ go install tool
```

### Run the Server and REST Proxy

```bash
$ go run cmd/server/main.go # starts gRPC server
$ go run cmd/rest/main.go # starts gRPC REST proxy
```

### View Swagger Documentation

```bash
$ make swagger-ui # visit http://localhost:80 to view Swagger UI
```

### Run Tests
```bash
$ make test # runs all tests
$ go test -v -run TestSimpleWires ./... # runs tests with a prefix
```

## TODO List

- [ ] Implement more bomb modules (Keypads, Button, Morse Code, etc.)
- [ ] Add config for game settings
  - [ ] Difficulty levels
  - [ ] Seeds for module generation
- [ ] Implement bomb timer and strike system
- [ ] Create comprehensive test suite
- [ ] Document gRPC API for client developers
- [x] Create a simple demo client ([In Progress](https://github.com/ZaneH/keep-talking-app))
- [x] Add CI/CD pipeline

### Bugs

- [ ] Bombs should never be entirely needy modules
- [ ] Maze Module: Player should never spawn on goal
- [ ] Wires can throw when index is OOB

## Contributing

Contributions are welcome. If you'd like to add a new module, fix a bug, or improve the codebase, feel free to open a pull request.

## License

This project is open-source under the MIT license.

**Note:** This project is an unofficial implementation inspired by "Keep Talking and Nobody Explodes" from Steel Crate Games. It is not affiliated with, endorsed by, or connected to Steel Crate Games in any way. All trademarks, game mechanics, and concepts belong to their respective owners. This implementation is created for educational purposes to explore software design patterns and architecture.

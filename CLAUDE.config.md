# Plan: Flexible Bomb Configuration System

## Overview
Implement a configuration system that allows API consumers to create games via three modes:
1. **Level-based** (1-10): Auto-generates difficulty-appropriate bombs
2. **Preset missions**: Enum-based selection matching vanilla KTANE missions
3. **Full custom**: Complete control over all bomb parameters with validation

## User Requirements
- Random modules should be **section-based** (only modules unlocked up to that section)
- CreateGameResponse returns only session_id (use GetBombs RPC for full bomb details)
- Reasonable limits on custom inputs to prevent abuse
- Mission specs support "pick one" logic via PossibleTypes field

---

## Implementation Steps

### Step 1: Create Proto Messages
**File:** `proto/game_config.proto` (NEW)

```protobuf
enum Mission { ... }  // THE_FIRST_BOMB, DOUBLE_YOUR_MONEY, etc.
message LevelConfig { int32 level = 1; }
message PresetMissionConfig { Mission mission = 1; }
message ModuleSpec { ModuleType type = 1; int32 count = 2; }
message CustomBombConfig { timer_seconds, max_strikes, num_faces, rows, columns, modules[], ... }
message GameConfig { oneof { level, preset, custom }; string seed = 10; }
message GeneratedConfigInfo { ... }  // For returning config details
```

**File:** `proto/player.proto` (MODIFY)
- Add `optional game_config.GameConfig config = 1` to `CreateGameRequest`
- Add `game_config.GeneratedConfigInfo config_info = 2` to `CreateGameResponse`

### Step 2: Define Missions and Module Pools
**File:** `internal/domain/valueobject/mission.go` (NEW)

Define all vanilla missions with their exact specs:
- Section 1-7 mission definitions (THE_FIRST_BOMB through exotic missions)
- Each mission has: name, timer, strikes, explicit module list, faces, grid size
- **Section-based module pools**: Define which modules are available in each section

```go
var SectionModulePools = map[int][]ModuleType{
    1: {BigButtonModule, KeypadModule, WiresModule},
    2: {BigButtonModule, KeypadModule, WiresModule},  // Same as intro
    3: {BigButtonModule, KeypadModule, WiresModule, PasswordModule, MorseModule,
        WireSequenceModule, ComplicatedWiresModule},
    4: {/* Section 3 + needy modules */},
    5: {/* All modules */},
    // ...
}
```

### Step 3: Define Level Configurations
**File:** `internal/domain/valueobject/level_config.go` (NEW)

Map levels 1-10 to difficulty parameters:
| Level | Timer | Strikes | Faces | Modules | Needy Weight |
|-------|-------|---------|-------|---------|--------------|
| 1 | 5:00 | 3 | 1 | 3 | 0% |
| 5 | 4:00 | 2 | 2 | 8 | 10% |
| 10 | 2:00 | 1 | 4 | 20 | 25% |

### Step 4: Add Validation
**File:** `internal/domain/valueobject/config_validation.go` (NEW)

Validation limits:
| Field | Min | Max |
|-------|-----|-----|
| timer_seconds | 30 | 3600 |
| max_strikes | 1 | 10 |
| num_faces | 1 | 10 |
| rows | 1 | 4 |
| columns | 1 | 5 |
| min_modules | 1 | 60 |
| max_modules_per_face | 1 | 15 |
| batteries | 0 | 6 |
| indicators | 0 | 5 |
| ports | 0 | 6 |

### Step 5: Update BombConfig
**File:** `internal/domain/valueobject/bomb_config.go` (MODIFY)

Add field for explicit module lists:
```go
ExplicitModules []ModuleSpec  // For missions with specific module requirements
```

### Step 6: Add Config Builder
**File:** `internal/domain/valueobject/config_builder.go` (NEW)

Functions to convert each config type to BombConfig:
- `FromLevel(level int) BombConfig`
- `FromMission(mission Mission) BombConfig`
- `FromCustom(proto *CustomBombConfig) BombConfig`

### Step 7: Update CreateGameCommand
**File:** `internal/application/command/create_game_command.go` (MODIFY)

```go
type CreateGameCommand struct {
    Seed         string
    ConfigType   ConfigType  // Default, Level, Mission, Custom
    Level        int
    Mission      Mission
    CustomConfig *BombConfig
}
```

### Step 8: Update GameService
**File:** `internal/application/services/game_service.go` (MODIFY)

Handle different config types in `CreateGameSession`:
```go
switch cmd.ConfigType {
case ConfigTypeLevel: config = NewGameSessionConfigFromLevel(...)
case ConfigTypeMission: config = NewGameSessionConfigFromMission(...)
case ConfigTypeCustom: config = NewGameSessionConfigFromCustom(...)
default: config = NewEasyGameSessionConfig(...)
}
```

Return generated config info for response.

### Step 9: Update BombFactory
**File:** `internal/domain/services/bomb_factory.go` (MODIFY)

Support explicit module lists:
```go
if len(config.ExplicitModules) > 0 {
    modules = expandExplicitModules(config.ExplicitModules)
} else {
    modules = generateWeightedModules(rng, config)
}
```

### Step 10: Update gRPC Adapter
**File:** `internal/infrastructure/grpc/game_service_adapter.go` (MODIFY)

- Convert proto config to domain command
- Return validation errors with proper gRPC status codes
- Populate GeneratedConfigInfo in response

### Step 11: Regenerate Proto
Run `make proto` or equivalent to regenerate Go bindings.

---

## File Summary

| File | Action |
|------|--------|
| `proto/game_config.proto` | CREATE |
| `proto/player.proto` | MODIFY |
| `internal/domain/valueobject/mission.go` | CREATE |
| `internal/domain/valueobject/level_config.go` | CREATE |
| `internal/domain/valueobject/config_validation.go` | CREATE |
| `internal/domain/valueobject/config_builder.go` | CREATE |
| `internal/domain/valueobject/bomb_config.go` | MODIFY |
| `internal/domain/valueobject/game_session_config.go` | MODIFY |
| `internal/application/command/create_game_command.go` | MODIFY |
| `internal/application/services/game_service.go` | MODIFY |
| `internal/domain/services/bomb_factory.go` | MODIFY |
| `internal/infrastructure/grpc/game_service_adapter.go` | MODIFY |

---

## API Examples

**Level-based:**
```json
POST /v1/game/create
{ "config": { "level": { "level": 5 } } }
```

**Preset mission:**
```json
POST /v1/game/create
{ "config": { "preset": { "mission": "BLINKENLIGHTS" } } }
```

**Full custom:**
```json
POST /v1/game/create
{
  "config": {
    "custom": {
      "timer_seconds": 180,
      "max_strikes": 1,
      "num_faces": 2,
      "rows": 2,
      "columns": 3,
      "modules": [
        { "type": "SIMON", "count": 3 },
        { "type": "MAZE", "count": 2 }
      ]
    }
  }
}
```

**Response:**
```json
{
  "session_id": "uuid"
}
```

To get bomb details, use `GetBombs` RPC with the session_id.

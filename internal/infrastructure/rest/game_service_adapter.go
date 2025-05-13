package rest

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/ZaneH/keep-talking/internal/application/command"
	"github.com/ZaneH/keep-talking/internal/application/services"
	"github.com/ZaneH/keep-talking/internal/domain/valueobject"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

// GameServiceAdapter adapts the application service to REST endpoints
type GameServiceAdapter struct {
	gameService *services.GameService
	router      *mux.Router
}

// NewGameServiceAdapter creates a new REST adapter for the game service
func NewGameServiceAdapter(gameService *services.GameService) *GameServiceAdapter {
	adapter := &GameServiceAdapter{
		gameService: gameService,
		router:      mux.NewRouter(),
	}
	adapter.setupRoutes()
	return adapter
}

// Router returns the configured router
func (a *GameServiceAdapter) Router() *mux.Router {
	return a.router
}

// setupRoutes configures all the REST endpoints
func (a *GameServiceAdapter) setupRoutes() {
	a.router.HandleFunc("/api/games", a.createGame).Methods("POST")
	a.router.HandleFunc("/api/games/{sessionId}/bombs", a.getBombs).Methods("GET")
	a.router.HandleFunc("/api/games/{sessionId}/input", a.sendInput).Methods("POST")
}

// createGame handles the creation of a new game session
func (a *GameServiceAdapter) createGame(w http.ResponseWriter, r *http.Request) {
	createGameCmd := &command.CreateGameCommand{}

	session, err := a.gameService.CreateGameSession(r.Context(), createGameCmd)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to create game: %v", err), http.StatusInternalServerError)
		return
	}

	log.Printf("Created game session with ID: %s\n", session.GetSessionID())

	response := struct {
		SessionID string `json:"sessionId"`
	}{
		SessionID: session.GetSessionID().String(),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

// getBombs retrieves all bombs for a given session
func (a *GameServiceAdapter) getBombs(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	sessionIDStr, ok := vars["sessionId"]
	if !ok {
		http.Error(w, "Session ID is required", http.StatusBadRequest)
		return
	}

	sessionID, err := uuid.Parse(sessionIDStr)
	if err != nil {
		http.Error(w, "Invalid session ID format", http.StatusBadRequest)
		return
	}

	gameState, err := a.gameService.GetGameSession(r.Context(), sessionID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get game session: %v", err), http.StatusInternalServerError)
		return
	}

	// Convert game state to a response format
	// This is a simplified version - you might want to create a more detailed response
	response := struct {
		SessionID string                 `json:"sessionId"`
		Bombs     map[string]interface{} `json:"bombs"`
	}{
		SessionID: sessionID.String(),
		Bombs:     make(map[string]interface{}),
	}

	for bombID, bombActor := range gameState.GetBombActors() {
		bomb := bombActor.GetBomb()
		bombData := map[string]interface{}{
			"id":           bomb.ID.String(),
			"serialNumber": bomb.SerialNumber,
			"strikeCount":  bomb.StrikeCount,
			"maxStrikes":   bomb.MaxStrikes,
			"batteries":    bomb.Batteries,
			"modules":      make(map[string]interface{}),
		}

		for moduleID, module := range bomb.Modules {
			moduleData := map[string]interface{}{
				"id":     moduleID.String(),
				"type":   module.GetType(),
				"solved": module.GetModuleState().MarkSolved,
			}
			bombData["modules"].(map[string]interface{})[moduleID.String()] = moduleData
		}

		response.Bombs[bombID.String()] = bombData
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// InputRequest represents the JSON structure for module input
type InputRequest struct {
	SessionID string      `json:"sessionId"`
	BombID    string      `json:"bombId"`
	ModuleID  string      `json:"moduleId"`
	InputType string      `json:"inputType"`
	Data      interface{} `json:"data"`
}

// sendInput handles player input for modules
func (a *GameServiceAdapter) sendInput(w http.ResponseWriter, r *http.Request) {
	var inputReq InputRequest
	if err := json.NewDecoder(r.Body).Decode(&inputReq); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	sessionID, err := uuid.Parse(inputReq.SessionID)
	if err != nil {
		http.Error(w, "Invalid session ID format", http.StatusBadRequest)
		return
	}

	bombID, err := uuid.Parse(inputReq.BombID)
	if err != nil {
		http.Error(w, "Invalid bomb ID format", http.StatusBadRequest)
		return
	}

	moduleID, err := uuid.Parse(inputReq.ModuleID)
	if err != nil {
		http.Error(w, "Invalid module ID format", http.StatusBadRequest)
		return
	}

	var cmd command.ModuleInputCommand

	switch inputReq.InputType {
	case "simpleWires":
		data, ok := inputReq.Data.(map[string]interface{})
		if !ok {
			http.Error(w, "Invalid data format for simple wires input", http.StatusBadRequest)
			return
		}

		wireIndex, ok := data["wireIndex"].(float64)
		if !ok {
			http.Error(w, "Wire index must be a number", http.StatusBadRequest)
			return
		}

		cmd = &command.SimpleWiresInputCommand{
			BaseModuleInputCommand: command.BaseModuleInputCommand{
				SessionID: sessionID,
				BombID:    bombID,
				ModuleID:  moduleID,
			},
			WireIndex: int(wireIndex),
		}

	case "password":
		data, ok := inputReq.Data.(map[string]interface{})
		if !ok {
			http.Error(w, "Invalid data format for password input", http.StatusBadRequest)
			return
		}

		action, ok := data["action"].(string)
		if !ok {
			http.Error(w, "Action field is required for password input", http.StatusBadRequest)
			return
		}

		if action == "letterChange" {
			letterIndex, ok := data["letterIndex"].(float64)
			if !ok {
				http.Error(w, "Letter index must be a number", http.StatusBadRequest)
				return
			}

			directionStr, ok := data["direction"].(string)
			if !ok {
				http.Error(w, "Direction must be a string", http.StatusBadRequest)
				return
			}

			var direction valueobject.IncrementDecrement
			if directionStr == "increment" {
				direction = valueobject.Increment
			} else if directionStr == "decrement" {
				direction = valueobject.Decrement
			} else {
				http.Error(w, "Direction must be 'increment' or 'decrement'", http.StatusBadRequest)
				return
			}

			cmd = &command.PasswordLetterChangeCommand{
				BaseModuleInputCommand: command.BaseModuleInputCommand{
					SessionID: sessionID,
					BombID:    bombID,
					ModuleID:  moduleID,
				},
				LetterIndex: int(letterIndex),
				Direction:   direction,
			}
		} else if action == "submit" {
			cmd = &command.PasswordSubmitCommand{
				BaseModuleInputCommand: command.BaseModuleInputCommand{
					SessionID: sessionID,
					BombID:    bombID,
					ModuleID:  moduleID,
				},
			}
		} else {
			http.Error(w, "Invalid action for password input", http.StatusBadRequest)
			return
		}

	case "bigButton":
		data, ok := inputReq.Data.(map[string]interface{})
		if !ok {
			http.Error(w, "Invalid data format for big button input", http.StatusBadRequest)
			return
		}

		pressTypeStr, ok := data["pressType"].(string)
		if !ok {
			http.Error(w, "Press type must be a string", http.StatusBadRequest)
			return
		}

		cmd = &command.BigButtonInputCommand{
			BaseModuleInputCommand: command.BaseModuleInputCommand{
				SessionID: sessionID,
				BombID:    bombID,
				ModuleID:  moduleID,
			},
			PressType: valueobject.PressType(pressTypeStr),
		}

	case "simonSays":
		data, ok := inputReq.Data.(map[string]interface{})
		if !ok {
			http.Error(w, "Invalid data format for Simon Says input", http.StatusBadRequest)
			return
		}

		colorStr, ok := data["color"].(string)
		if !ok {
			http.Error(w, "Color must be a string", http.StatusBadRequest)
			return
		}

		cmd = &command.SimonSaysInputCommand{
			BaseModuleInputCommand: command.BaseModuleInputCommand{
				SessionID: sessionID,
				BombID:    bombID,
				ModuleID:  moduleID,
			},
			Color: valueobject.Color(colorStr),
		}

	default:
		http.Error(w, fmt.Sprintf("Unknown input type: %s", inputReq.InputType), http.StatusBadRequest)
		return
	}

	res, err := a.gameService.ProcessModuleInput(r.Context(), cmd)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to process input: %v", err), http.StatusInternalServerError)
		return
	}

	// Create a response based on the command result
	response := struct {
		Success bool        `json:"success"`
		Solved  bool        `json:"solved"`
		Strike  bool        `json:"strike"`
		Data    interface{} `json:"data,omitempty"`
	}{
		Success: true,
	}

	// Handle different result types
	switch result := res.(type) {
	case *command.SimpleWiresInputCommandResult:
		response.Solved = result.Solved
		response.Strike = result.Strike

	case *command.PasswordLetterChangeCommandResult:
		response.Solved = result.Solved
		response.Strike = result.Strike

	case *command.PasswordSubmitCommandResult:
		response.Solved = result.Solved
		response.Strike = result.Strike

	case *command.BigButtonInputCommandResult:
		response.Solved = result.Solved
		response.Strike = result.Strike
		if result.StripColor != "" {
			response.Data = map[string]string{
				"stripColor": string(result.StripColor),
			}
		}

	case *command.SimonSaysInputCommandResult:
		response.Solved = result.Solved
		response.Strike = result.Strike
		if len(result.NextColorSequence) > 0 {
			colorStrings := make([]string, len(result.NextColorSequence))
			for i, color := range result.NextColorSequence {
				colorStrings[i] = string(color)
			}
			response.Data = map[string]interface{}{
				"nextColorSequence": colorStrings,
			}
		}

	default:
		log.Printf("Unknown result type: %T", res)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

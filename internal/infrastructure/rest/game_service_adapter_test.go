package rest

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ZaneH/keep-talking/internal/actors"
	"github.com/ZaneH/keep-talking/internal/application/services"
	domainservices "github.com/ZaneH/keep-talking/internal/domain/services"
	"github.com/ZaneH/keep-talking/internal/infrastructure/adapters"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func setupTestAdapter() *GameServiceAdapter {
	actorSystem := actors.NewActorSystem()
	actorSystemAdapter := adapters.NewActorSystemAdapter(actorSystem)
	bombFactory := &domainservices.BombFactoryImpl{}
	bombService := services.NewBombService(actorSystemAdapter, bombFactory)
	gameService := services.NewGameService(actorSystem, bombService)

	return NewGameServiceAdapter(gameService)
}

func TestCreateGame(t *testing.T) {
	// Arrange
	adapter := setupTestAdapter()
	req, err := http.NewRequest("POST", "/api/games", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(adapter.createGame)

	// Act
	handler.ServeHTTP(rr, req)

	// Assert
	assert.Equal(t, http.StatusCreated, rr.Code)

	var response struct {
		SessionID string `json:"sessionId"`
	}
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.NotEmpty(t, response.SessionID)

	// Verify the session ID is a valid UUID
	_, err = uuid.Parse(response.SessionID)
	assert.NoError(t, err)
}

func TestSendInput(t *testing.T) {
	// This is a more complex test that would require setting up a game session
	// with bombs and modules. For simplicity, we'll just test the API structure.

	// Arrange
	adapter := setupTestAdapter()

	// Create a game session first
	createReq, _ := http.NewRequest("POST", "/api/games", nil)
	createRR := httptest.NewRecorder()
	http.HandlerFunc(adapter.createGame).ServeHTTP(createRR, createReq)

	var createResp struct {
		SessionID string `json:"sessionId"`
	}
	json.Unmarshal(createRR.Body.Bytes(), &createResp)

	// For a full test, we would need to add a bomb with modules to this session
	// Since we can't easily do that in this test, we'll just verify the API rejects
	// invalid requests properly

	inputReq := InputRequest{
		SessionID: createResp.SessionID,
		BombID:    uuid.New().String(), // Random bomb ID that doesn't exist
		ModuleID:  uuid.New().String(), // Random module ID that doesn't exist
		InputType: "simpleWires",
		Data: map[string]interface{}{
			"wireIndex": 2,
		},
	}

	body, _ := json.Marshal(inputReq)
	req, _ := http.NewRequest("POST", "/api/games/"+createResp.SessionID+"/input", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	http.HandlerFunc(adapter.sendInput).ServeHTTP(rr, req)

	// We expect an error since the bomb doesn't exist
	assert.NotEqual(t, http.StatusOK, rr.Code)
}

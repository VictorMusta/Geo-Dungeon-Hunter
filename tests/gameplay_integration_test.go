package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"dungeons/app/models"
	"github.com/stretchr/testify/assert"
)

func TestFullGameplayCycle(t *testing.T) {
	db := SetupTestDB(t)
	r := SetupFullTestRouter(db)

	// 1. Create Player
	playerData := models.Player{
		DisplayName: "Hero",
		Password:    "password123",
	}
	body, _ := json.Marshal(playerData)
	req := httptest.NewRequest("POST", "/v1/players", bytes.NewBuffer(body))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code, "Create player")

	// 2. Login - use 'display_name' matching the LoginRequest model JSON tag
	loginData := map[string]string{
		"display_name": "Hero",
		"password":     "password123",
	}
	body, _ = json.Marshal(loginData)
	req = httptest.NewRequest("POST", "/v1/players/login", bytes.NewBuffer(body))
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code, "Login")
	t.Logf("Login resp: %s", w.Body.String())

	var loginResp models.LoginResponse
	json.Unmarshal(w.Body.Bytes(), &loginResp)
	token := loginResp.Token
	playerID := loginResp.Player.CustomID
	assert.NotEmpty(t, token, "Expected non-empty token after login")
	assert.NotEmpty(t, playerID, "Expected non-empty playerID after login")

	// 3. Create Dungeon - GM route is /v1/mj/dungeons; Area.Name is required
	dungeonData := models.Dungeon{
		Title:       "Dark Cave",
		Description: "A scary cave",
		CreatedBy:   playerID,
		Area:        models.Area{Name: "Forest Zone"},
	}
	body, _ = json.Marshal(dungeonData)
	req = httptest.NewRequest("POST", "/v1/mj/dungeons", bytes.NewBuffer(body))
	req.Header.Set("Authorization", "Bearer "+token)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	t.Logf("Create Dungeon resp: %s", w.Body.String())
	assert.Equal(t, http.StatusCreated, w.Code, "Create dungeon")

	// Dungeon is returned inside WSResponse.Data
	var dungeonWrap map[string]json.RawMessage
	json.Unmarshal(w.Body.Bytes(), &dungeonWrap)
	var dungeonResp models.Dungeon
	json.Unmarshal(dungeonWrap["data"], &dungeonResp)
	dungeonID := dungeonResp.CustomID
	assert.NotEmpty(t, dungeonID, "Expected non-empty dungeonID")

	// 4. Add Boss Step
	stepData := models.BossStep{
		Name:       "Giant Spider",
		GoldReward: 100,
		Difficulty: 5,
		Location:   models.Location{Lat: 48.8566, Lon: 2.3522, RadiusMeters: 100},
	}
	body, _ = json.Marshal(stepData)
	req = httptest.NewRequest("POST", fmt.Sprintf("/v1/mj/dungeons/%s/steps", dungeonID), bytes.NewBuffer(body))
	req.Header.Set("Authorization", "Bearer "+token)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	t.Logf("Create Step resp: %s", w.Body.String())
	assert.Equal(t, http.StatusCreated, w.Code, "Create step")

	var stepWrap map[string]json.RawMessage
	json.Unmarshal(w.Body.Bytes(), &stepWrap)
	var stepResp models.BossStep
	json.Unmarshal(stepWrap["data"], &stepResp)
	stepID := stepResp.CustomID
	assert.NotEmpty(t, stepID, "Expected non-empty stepID")

	// 5. Publish Dungeon
	req = httptest.NewRequest("POST", fmt.Sprintf("/v1/mj/dungeons/%s/publish", dungeonID), nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	t.Logf("Publish resp: %s", w.Body.String())
	assert.Equal(t, http.StatusOK, w.Code, "Publish dungeon")

	// 6. Start Run
	runData := models.Run{
		DungeonID: dungeonID,
		PlayerID:  playerID,
	}
	body, _ = json.Marshal(runData)
	req = httptest.NewRequest("POST", "/v1/runs", bytes.NewBuffer(body))
	req.Header.Set("Authorization", "Bearer "+token)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	t.Logf("Create Run resp: %s", w.Body.String())
	assert.Equal(t, http.StatusCreated, w.Code, "Create run")

	var runWrap map[string]json.RawMessage
	json.Unmarshal(w.Body.Bytes(), &runWrap)
	var runResp models.Run
	json.Unmarshal(runWrap["data"], &runResp)
	runID := runResp.CustomID
	assert.NotEmpty(t, runID, "Expected non-empty runID")

	// 7. Attempt Boss (In Range)
	attemptData := map[string]float64{
		"lat": 48.8566,
		"lon": 2.3522,
	}
	body, _ = json.Marshal(attemptData)
	req = httptest.NewRequest("POST", fmt.Sprintf("/v1/runs/%s/steps/%s/attempt", runID, stepID), bytes.NewBuffer(body))
	req.Header.Set("Authorization", "Bearer "+token)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	t.Logf("Attempt resp: %s", w.Body.String())
	assert.Equal(t, http.StatusOK, w.Code, "Boss attempt")

	// 8. Verify Rewards (Gold) - GET /players/:id is protected, needs token
	req = httptest.NewRequest("GET", "/v1/players/"+playerID, nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code, "Get player by ID")

	var playerWrap map[string]json.RawMessage
	json.Unmarshal(w.Body.Bytes(), &playerWrap)
	var playerResp models.Player
	json.Unmarshal(playerWrap["data"], &playerResp)
	t.Logf("Player Gold: %d (raw: %s)", playerResp.Gold, w.Body.String())
	assert.Equal(t, int64(100), playerResp.Gold)
}

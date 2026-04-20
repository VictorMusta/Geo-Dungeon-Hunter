package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPlayerRegistrationAndLogin(t *testing.T) {
	db := SetupTestDB(t)
	router := SetupFullTestRouter(db)

	playerName := "TestWarrior"
	password := "password123"

	// 1. Test Registration
	regPayload := map[string]interface{}{
		"display_name": playerName,
		"password":     password,
	}
	body, _ := json.Marshal(regPayload)
	req, _ := http.NewRequest("POST", "/v1/players", bytes.NewBuffer(body))
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("Expected 201, got %d. Body: %s", w.Code, w.Body.String())
	}

	var regResponse map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &regResponse)
	
	token := regResponse["token"].(string)
	if token == "" {
		t.Error("Expected token in registration response, got empty")
	}

	// 2. Test Login
	loginPayload := map[string]interface{}{
		"display_name": playerName,
		"password":     password,
	}
	body, _ = json.Marshal(loginPayload)
	req, _ = http.NewRequest("POST", "/v1/players/login", bytes.NewBuffer(body))
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected 200 on login, got %d. Body: %s", w.Code, w.Body.String())
	}

	var loginResponse map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &loginResponse)
	
	if loginResponse["token"] == "" {
		t.Error("Expected token in login response, got empty")
	}
}

func TestLogin_IncorrectPassword(t *testing.T) {
	db := SetupTestDB(t)
	router := SetupFullTestRouter(db)

	playerName := "TestWarrior"
	
	// Create player first
	regPayload := map[string]interface{}{
		"display_name": playerName,
		"password":     "correct_password",
	}
	body, _ := json.Marshal(regPayload)
	req, _ := http.NewRequest("POST", "/v1/players", bytes.NewBuffer(body))
	router.ServeHTTP(httptest.NewRecorder(), req)

	// Try login with wrong password
	loginPayload := map[string]interface{}{
		"display_name": playerName,
		"password":     "wrong_password",
	}
	body, _ = json.Marshal(loginPayload)
	req, _ = http.NewRequest("POST", "/v1/players/login", bytes.NewBuffer(body))
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("Expected 401 for wrong password, got %d", w.Code)
	}
}

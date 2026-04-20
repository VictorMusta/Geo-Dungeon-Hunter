package functions

import (
	"dungeons/app/server"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// GenerateToken creates a new JWT for a player
func GenerateToken(playerID string) (string, error) {
	srv := server.GetServer()
	
	claims := jwt.MapClaims{
		"sub": playerID,
		"exp": time.Now().Add(time.Hour * 24).Unix(), // valid for 24 hours
		"iat": time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(srv.TokenKey))
}

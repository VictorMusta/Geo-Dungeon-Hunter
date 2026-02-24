package server

import (
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

var server *Dungeons

// Dungeons Structure
type Dungeons struct {
	Database  *mongo.Database
	Router    *gin.Engine
	Version   string
	Port      string
	TokenKey  string
	Origin    string
	LogFormat string
	Mode      string
	DBHost    string
}

func (d *Dungeons) ParseParameters() {
	d.LogFormat = os.Getenv("LOG_FORMAT")
	d.Version = os.Getenv("API_VERSION")
	d.Port = os.Getenv("API_PORT")
	d.TokenKey = os.Getenv("TOKEN_KEY")
	d.Origin = os.Getenv("ALLOW_ORIGIN")
	d.Mode = os.Getenv("MODE")
	d.DBHost = os.Getenv("DB_HOST")
}

// ListenAndServe listens on the TCP network address addr and then calls Serve with handler to handle requests on incoming connections.
// https://github.com/gin-gonic/gin
func (d *Dungeons) ListenAndServe() error {
	srv := &http.Server{
		Addr:              d.Port,
		Handler:           d.Router,
		ReadHeaderTimeout: 10 * time.Second,
	}

	// start
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal().Msgf("Unable to listen and serve: %v", err)
		return err
	}
	return nil
}

// SetServer init mongo database
func SetServer(s *Dungeons) {
	server = s
}

// GetServer Flashcards
func GetServer() *Dungeons {
	return server
}

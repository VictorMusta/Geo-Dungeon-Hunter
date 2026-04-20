package tests

import (
	"context"
	"dungeons/app/mongodb"
	"os"
	"testing"
	"dungeons/app/server"
	
	"dungeons/app/routes/auction"
	"dungeons/app/routes/dungeon"
	"dungeons/app/routes/inventory"
	"dungeons/app/routes/item"
	"dungeons/app/routes/leaderboard"
	"dungeons/app/routes/player"
	"dungeons/app/routes/run"
	"dungeons/app/routes/common"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

// SetupTestDB initializes a connection to the test database
func SetupTestDB(t *testing.T) *mongo.Database {
	t.Helper()

	dbHost := os.Getenv("DB_HOST")
	if dbHost == "" {
		dbHost = "mongodb://localhost:27017"
	}

	client, err := mongodb.OpenMongoDB(dbHost)
	if err != nil {
		t.Fatalf("Failed to connect to test DB: %v", err)
	}

	db := client.Database("dungeons_test")
	
	// Ensure DB is clean
	err = db.Drop(context.TODO())
	if err != nil {
		t.Fatalf("Failed to drop test DB: %v", err)
	}

	return db
}

// SetupFullTestRouter initializes a router with all domain routes wired for hexagonal DI
func SetupFullTestRouter(db *mongo.Database) *gin.Engine {
	server.SetServer(&server.Dungeons{
		Database: db,
		TokenKey: "test-secret-key",
	})

	gin.SetMode(gin.TestMode)
	r := common.SetupRouter()
	
	player.SetupRouter(r)
	item.SetupRouter(r)
	dungeon.SetupRouter(r)
	run.SetupRouter(r)
	inventory.SetupRouter(r)
	auction.SetupRouter(r)
	leaderboard.SetupRouter(r)

	return r
}

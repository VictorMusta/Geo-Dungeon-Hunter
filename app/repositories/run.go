package repositories

import (
	"context"
	"dungeons/app/models"
)

type RunRepository interface {
	GetByPlayerID(playerID string) ([]models.Run, error)
	GetByID(id string) (models.Run, error)
	GetActiveRun(playerID, dungeonId string) (models.Run, error)
	Create(run *models.Run) error
	Update(id string, run *models.Run) error
	
	// Transactional operations for gameplay
	ExecuteBossAttempt(ctx context.Context, runID string, playerID string, rewards models.RewardsGiven, isCompleted bool, killedStep models.KilledStep) error
}

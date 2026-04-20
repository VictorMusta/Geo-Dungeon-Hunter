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
	ExecuteBossAttempt(ctx context.Context, runID string, killedStep models.KilledStep) error
	CommitRewards(ctx context.Context, runID string, playerID string, gold int64, items []models.RewardItem, newState string) error
}

package repositories

import (
	"dungeons/app/models"
)

// PlayerRepository defines the port for player data access
type PlayerRepository interface {
	Get(params models.QueryParams) ([]models.Player, error)
	GetByID(id string) (models.Player, error)
	GetByIDs(ids []string) ([]models.Player, error)
	Create(player *models.Player) error
	Update(id string, player *models.Player) error
	Suspend(id string) error
	FindByDisplayName(name string) (models.Player, error)
}

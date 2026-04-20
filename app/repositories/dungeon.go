package repositories

import (
	"dungeons/app/models"
)

type DungeonRepository interface {
	Get(params models.QueryParams) ([]models.Dungeon, error)
	GetByID(id string) (models.Dungeon, error)
	Create(dungeon *models.Dungeon) error
	Update(id string, dungeon *models.Dungeon) error
	GetPublished(params models.QueryParams) ([]models.Dungeon, error)
	GetByMJ(mjId string, params models.QueryParams) ([]models.Dungeon, error)
}

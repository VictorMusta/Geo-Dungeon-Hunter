package repositories

import (
	"dungeons/app/models"
)

type BossStepRepository interface {
	GetByDungeon(dungeonId string) ([]models.BossStep, error)
	GetByDungeonOrdered(dungeonId string) ([]models.BossStep, error)
	GetByID(dungeonId, id string) (models.BossStep, error)
	Create(step *models.BossStep) error
	Update(dungeonId, id string, step *models.BossStep) error
	Delete(dungeonId, id string) error
	CountByDungeon(dungeonId string) (int64, error)
}

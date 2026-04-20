package repositories

import (
	"dungeons/app/models"
)

type ItemRepository interface {
	Get(params models.QueryParams) ([]models.ItemDef, error)
	GetByID(id string) (models.ItemDef, error)
	Create(item *models.ItemDef) error
	Update(id string, item *models.ItemDef) error
}

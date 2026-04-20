package repositories

import (
	"dungeons/app/models"
)

type InventoryRepository interface {
	GetByPlayerID(playerID string) ([]models.InventoryEntry, error)
	GetByItem(playerID, itemID string) (models.InventoryEntry, error)
	Update(playerID string, itemID string, qtyDelta int64) error
}

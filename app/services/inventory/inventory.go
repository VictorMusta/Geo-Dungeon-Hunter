package inventory

import (
	"dungeons/app/models"
	"dungeons/app/repositories"
)

type Inventory struct {
	repo repositories.InventoryRepository
}

func New(repo repositories.InventoryRepository) *Inventory {
	return &Inventory{repo: repo}
}

func (s *Inventory) GetByPlayerID(playerID string) (*models.InventoryResponse, error) {
	entries, err := s.repo.GetByPlayerID(playerID)
	if err != nil {
		return nil, err
	}

	var items []models.InventoryItemDTO
	for _, entry := range entries {
		items = append(items, models.InventoryItemDTO{
			ItemID: entry.ItemID,
			Qty:    entry.Qty,
		})
	}

	if items == nil {
		items = []models.InventoryItemDTO{}
	}

	return &models.InventoryResponse{
		PlayerID: playerID,
		Items:    items,
	}, nil
}

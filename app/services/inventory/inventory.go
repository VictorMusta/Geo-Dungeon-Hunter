package inventory

import (
	"dungeons/app/models"
	"dungeons/app/repositories"
)

type Inventory struct {
	repo     repositories.InventoryRepository
	itemRepo repositories.ItemRepository
}

func New(repo repositories.InventoryRepository, itemRepo repositories.ItemRepository) *Inventory {
	return &Inventory{repo: repo, itemRepo: itemRepo}
}

func (s *Inventory) GetByPlayerID(playerID string) (*models.InventoryResponse, error) {
	entries, err := s.repo.GetByPlayerID(playerID)
	if err != nil {
		return nil, err
	}

	var items []models.InventoryItemDTO
	for _, entry := range entries {
		dto := models.InventoryItemDTO{
			ItemID: entry.ItemID,
			Qty:    entry.Qty,
		}

		// Fetch item details
		item, err := s.itemRepo.GetByID(entry.ItemID)
		if err == nil {
			dto.Item = &models.ItemDefResponse{
				ID:          item.CustomID,
				Name:        item.Name,
				Type:        item.Type,
				Rarity:      item.Rarity,
				Description: item.Description,
				Tradable:    item.Tradable,
				BaseValue:   item.BaseValue,
				Stats:       item.Stats,
			}
		}

		items = append(items, dto)
	}

	if items == nil {
		items = []models.InventoryItemDTO{}
	}

	return &models.InventoryResponse{
		PlayerID: playerID,
		Items:    items,
	}, nil
}

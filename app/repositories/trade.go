package repositories

import (
	"dungeons/app/models"
)

type TradeRepository interface {
	Create(trade *models.Trade) error
	GetByPlayer(playerID string) ([]models.Trade, error)
}

package models

import "time"

type ItemID string

type InventoryEntry struct {
	PlayerID  string    `db:"player_id"`
	ItemID    string    `db:"item_id"`
	Qty       int64     `db:"qty"`
	UpdatedAt time.Time `db:"updated_at"`
}

// Collection Mongodb collection
func (ie *InventoryEntry) Collection() string {
	return "inventory"
}

type ItemDef struct {
	ID          string    `db:"id"`
	Type        string    `db:"type"`
	Rarity      string    `db:"rarity"`
	Name        string    `db:"name"`
	Description string    `db:"description"`
	StatsJSON   []byte    `db:"stats_json"` // JSONB côté Postgres
	Tradable    bool      `db:"tradable"`
	BaseValue   int64     `db:"base_value"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}

// Collection Mongodb collection
func (i *ItemDef) Collection() string {
	return "item"
}

type InventoryResponse struct {
	PlayerID string             `json:"playerId"`
	Items    []InventoryItemDTO `json:"items"`
}

type InventoryItemDTO struct {
	ItemID string `json:"itemId"`
	Qty    int64  `json:"qty"`
}

type ItemDefResponse struct {
	ID          string         `json:"id"`
	Type        string         `json:"type"`
	Rarity      string         `json:"rarity"`
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Stats       map[string]any `json:"stats,omitempty"`
	Tradable    bool           `json:"tradable"`
	BaseValue   int64          `json:"baseValue,omitempty"`
}

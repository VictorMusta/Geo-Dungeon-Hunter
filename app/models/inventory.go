package models

import "time"

type ItemID string

type InventoryEntry struct {
	PlayerID  string    `bson:"playerId" json:"playerId"`
	ItemID    string    `bson:"itemId" json:"itemId"`
	Qty       int64     `bson:"qty" json:"qty"`
	UpdatedAt time.Time `bson:"updatedAt" json:"updatedAt"`
}

func (ie *InventoryEntry) Collection() string {
	return "inventory"
}

type ItemDef struct {
	CustomID    string         `bson:"customID" json:"id"`
	Type        string         `bson:"type" json:"type" validate:"required,oneof=weapon artifact consumable"`
	Rarity      string         `bson:"rarity" json:"rarity" validate:"required,oneof=common uncommon rare epic legendary"`
	Name        string         `bson:"name" json:"name" validate:"required"`
	Description string         `bson:"description" json:"description"`
	Stats       map[string]any `bson:"stats,omitempty" json:"stats,omitempty"`
	Tradable    bool           `bson:"tradable" json:"tradable"`
	BaseValue   int64          `bson:"baseValue" json:"baseValue"`
	CreatedAt   time.Time      `bson:"createdAt" json:"createdAt"`
	UpdatedAt   time.Time      `bson:"updatedAt" json:"updatedAt"`
}

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

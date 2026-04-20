package models

import "time"

type Location struct {
	Lat          float64 `bson:"lat" json:"lat" validate:"required"`
	Lon          float64 `bson:"lon" json:"lon" validate:"required"`
	RadiusMeters float64 `bson:"radiusMeters" json:"radiusMeters" validate:"required,gt=0"`
}

type LootEntry struct {
	ItemID   string  `bson:"itemId" json:"itemId" validate:"required"`
	DropRate float64 `bson:"dropRate" json:"dropRate" validate:"required,gt=0,lte=1"`
	MinQty   int     `bson:"minQty" json:"minQty" validate:"required,min=1"`
	MaxQty   int     `bson:"maxQty" json:"maxQty" validate:"required,gtefield=MinQty"`
}

type BossStep struct {
	CustomID        string      `bson:"customID" json:"id"`
	DungeonID       string      `bson:"dungeonId" json:"dungeonId" validate:"required"`
	Order           int         `bson:"order" json:"order"`
	Name            string      `bson:"name" json:"name" validate:"required"`
	Emoji           string      `bson:"emoji" json:"emoji"`
	Location        Location    `bson:"location" json:"location" validate:"required"`
	ZoneDescription string      `bson:"zoneDescription" json:"zoneDescription"`
	Difficulty      int         `bson:"difficulty" json:"difficulty" validate:"required,min=1,max=10"`
	GoldReward      int64       `bson:"goldReward" json:"goldReward"`
	LootTable       []LootEntry `bson:"lootTable,omitempty" json:"lootTable,omitempty"`
	CreatedAt       time.Time   `bson:"createdAt" json:"createdAt"`
	UpdatedAt       time.Time   `bson:"updatedAt" json:"updatedAt"`
}

func (b *BossStep) Collection() string {
	return "boss_step"
}

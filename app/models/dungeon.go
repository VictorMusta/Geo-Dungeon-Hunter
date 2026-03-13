package models

import "time"

type BoundingBox struct {
	MinLat float64 `bson:"minLat" json:"minLat"`
	MaxLat float64 `bson:"maxLat" json:"maxLat"`
	MinLon float64 `bson:"minLon" json:"minLon"`
	MaxLon float64 `bson:"maxLon" json:"maxLon"`
}

type Area struct {
	Name        string       `bson:"name" json:"name" validate:"required"`
	BoundingBox *BoundingBox `bson:"boundingBox,omitempty" json:"boundingBox,omitempty"`
}

type Dungeon struct {
	CustomID    string    `bson:"customID" json:"id"`
	Title       string    `bson:"title" json:"title" validate:"required,min=3"`
	Description string    `bson:"description" json:"description"`
	CreatedBy   string    `bson:"createdBy" json:"createdBy" validate:"required"`
	Area        Area      `bson:"area" json:"area"`
	Status      string    `bson:"status" json:"status"`
	CreatedAt   time.Time `bson:"createdAt" json:"createdAt"`
	UpdatedAt   time.Time `bson:"updatedAt" json:"updatedAt"`
}

func (d *Dungeon) Collection() string {
	return "dungeon"
}

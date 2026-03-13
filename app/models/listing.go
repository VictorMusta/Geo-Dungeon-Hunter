package models

import "time"

type Listing struct {
	CustomID     string    `bson:"customID" json:"id"`
	SellerID     string    `bson:"sellerId" json:"sellerId" validate:"required"`
	ItemID       string    `bson:"itemId" json:"itemId" validate:"required"`
	Qty          int       `bson:"qty" json:"qty" validate:"required,min=1"`
	PricePerUnit int64     `bson:"pricePerUnit" json:"pricePerUnit" validate:"required,min=1"`
	Status       string    `bson:"status" json:"status"`
	CreatedAt    time.Time `bson:"createdAt" json:"createdAt"`
	ExpiresAt    time.Time `bson:"expiresAt" json:"expiresAt"`
}

func (l *Listing) Collection() string {
	return "listing"
}

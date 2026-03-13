package models

import "time"

type Trade struct {
	CustomID   string    `bson:"customID" json:"id"`
	BuyerID    string    `bson:"buyerId" json:"buyerId"`
	SellerID   string    `bson:"sellerId" json:"sellerId"`
	ListingID  string    `bson:"listingId" json:"listingId"`
	Qty        int       `bson:"qty" json:"qty"`
	TotalPrice int64     `bson:"totalPrice" json:"totalPrice"`
	CreatedAt  time.Time `bson:"createdAt" json:"createdAt"`
}

func (t *Trade) Collection() string {
	return "trade"
}

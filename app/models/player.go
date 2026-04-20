package models

import "time"

type PlayerID string

type Player struct {
	CustomID    string    `bson:"customID" json:"id"`
	DisplayName string    `bson:"display_name" json:"display_name" validate:"required,min=3"`
	Password    string    `bson:"password" json:"password,omitempty" validate:"omitempty,min=6"`
	Gold        int64     `bson:"gold" json:"gold"`
	CreatedAt   time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt   time.Time `bson:"updated_at" json:"updated_at"`
}

type PlayerResponse struct {
	ID          string    `json:"id"`
	DisplayName string    `json:"displayName"`
	Wallet      Wallet    `json:"wallet"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

type Wallet struct {
	Gold int64 `json:"gold"` // int64 pour éviter les soucis quand ça grossit
}

// Collection Mongodb collection
func (p *Player) Collection() string {
	return "player"
}

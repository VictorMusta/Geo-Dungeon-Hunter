package models

import "time"

type PlayerID string

type Player struct {
	CustomID    string    `bson:"customID" json:"id"`
	DisplayName string    `json:"display_name"`
	Gold        int64     `json:"gold"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
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

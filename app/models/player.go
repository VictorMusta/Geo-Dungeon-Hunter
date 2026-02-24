package models

import "time"

type PlayerID string

type Player struct {
	ID          string    `db:"id"`
	DisplayName string    `db:"display_name"`
	Gold        int64     `db:"gold"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
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

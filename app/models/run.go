package models

import "time"

type RewardItem struct {
	ItemID      string           `bson:"itemId" json:"itemId"`
	Qty         int              `bson:"qty" json:"qty"`
	ItemDetails *ItemDefResponse `bson:"itemDetails,omitempty" json:"itemDetails,omitempty"`
}

type RewardsGiven struct {
	Gold  int64        `bson:"gold" json:"gold"`
	Items []RewardItem `bson:"items,omitempty" json:"items,omitempty"`
}

type KilledStep struct {
	BossStepID   string       `bson:"bossStepId" json:"bossStepId"`
	KilledAt     time.Time    `bson:"killedAt" json:"killedAt"`
	RewardsGiven RewardsGiven `bson:"rewardsGiven" json:"rewardsGiven"`
}

type Run struct {
	CustomID    string       `bson:"customID" json:"id"`
	DungeonID   string       `bson:"dungeonId" json:"dungeonId" validate:"required"`
	PlayerID    string       `bson:"playerId" json:"playerId" validate:"required"`
	State       string       `bson:"state" json:"state"`
	CurrentStep int          `bson:"currentStep" json:"currentStep"`
	KilledSteps []KilledStep `bson:"killedSteps,omitempty" json:"killedSteps,omitempty"`
	StartedAt   time.Time    `bson:"startedAt" json:"startedAt"`
	EndedAt     *time.Time   `bson:"endedAt,omitempty" json:"endedAt,omitempty"`
	UpdatedAt   time.Time    `bson:"updatedAt" json:"updatedAt"`
}

func (r *Run) Collection() string {
	return "run"
}

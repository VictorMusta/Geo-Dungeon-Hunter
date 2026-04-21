package repositories

type LeaderboardEntry struct {
	PlayerID    string  `json:"playerId" bson:"playerId"`
	DisplayName string  `json:"displayName" bson:"displayName"`
	Score       float64 `json:"score" bson:"score"`
}

type LeaderboardRepository interface {
	GetByCompletions(limit int) ([]LeaderboardEntry, error)
	GetByGold(limit int) ([]LeaderboardEntry, error)
}

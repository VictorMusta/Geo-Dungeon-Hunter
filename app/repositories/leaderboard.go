package repositories

type LeaderboardEntry struct {
	PlayerID    string  `json:"playerId"`
	DisplayName string  `json:"displayName"`
	Score       float64 `json:"score"`
}

type LeaderboardRepository interface {
	GetByCompletions(limit int) ([]LeaderboardEntry, error)
	GetByGold(limit int) ([]LeaderboardEntry, error)
	GetBySpeed(dungeonID string, limit int) ([]LeaderboardEntry, error)
}

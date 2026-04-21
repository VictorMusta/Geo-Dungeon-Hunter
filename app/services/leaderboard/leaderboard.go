package leaderboard

import (
	"dungeons/app/repositories"
)

type Leaderboard struct {
	repo repositories.LeaderboardRepository
}

func New(repo repositories.LeaderboardRepository) *Leaderboard {
	return &Leaderboard{repo: repo}
}

func (s *Leaderboard) GetByCompletions(limit int) ([]repositories.LeaderboardEntry, error) {
	return s.repo.GetByCompletions(limit)
}

func (s *Leaderboard) GetByGold(limit int) ([]repositories.LeaderboardEntry, error) {
	return s.repo.GetByGold(limit)
}

package run

import (
	"context"
	"dungeons/app/functions"
	"dungeons/app/models"
	"dungeons/app/repositories"
	"errors"
	"math/rand"
	"time"

	"github.com/go-playground/validator/v10"
)

type Run struct {
	repo       repositories.RunRepository
	dungeonRepo repositories.DungeonRepository
	bsRepo     repositories.BossStepRepository
	playerRepo repositories.PlayerRepository
	validate   *validator.Validate
}

func New(
	repo repositories.RunRepository,
	dungeonRepo repositories.DungeonRepository,
	bsRepo repositories.BossStepRepository,
	playerRepo repositories.PlayerRepository,
) *Run {
	return &Run{
		repo:        repo,
		dungeonRepo: dungeonRepo,
		bsRepo:      bsRepo,
		playerRepo:  playerRepo,
		validate:    validator.New(),
	}
}

func (s *Run) Create(in *models.Run) (*models.Run, error) {
	// Verify dungeon exists and is published
	d, err := s.dungeonRepo.GetByID(in.DungeonID)
	if err != nil {
		return nil, errors.New("dungeon not found")
	}
	if d.Status != "published" {
		return nil, errors.New("dungeon is not published")
	}

	// Verify player exists
	_, err = s.playerRepo.GetByID(in.PlayerID)
	if err != nil {
		return nil, errors.New("player not found")
	}

	// Check no active run for this player+dungeon
	_, err = s.repo.GetActiveRun(in.PlayerID, in.DungeonID)
	if err == nil {
		return nil, errors.New("player already has an active run for this dungeon")
	}

	run := models.Run{
		CustomID:    functions.NewUUID(),
		DungeonID:   in.DungeonID,
		PlayerID:    in.PlayerID,
		State:       "active",
		CurrentStep: 1,
		KilledSteps: []models.KilledStep{},
		StartedAt:   time.Now(),
	}

	if err := s.repo.Create(&run); err != nil {
		return nil, err
	}

	return &run, nil
}

func (s *Run) GetByID(id string) (models.Run, error) {
	return s.repo.GetByID(id)
}

func (s *Run) GetByPlayerID(playerID string) ([]models.Run, error) {
	return s.repo.GetByPlayerID(playerID)
}

func (s *Run) Abandon(id string) error {
	r, err := s.repo.GetByID(id)
	if err != nil {
		return errors.New("run not found")
	}
	if r.State != "active" {
		return errors.New("only active runs can be abandoned")
	}

	now := time.Now()
	r.State = "abandoned"
	r.EndedAt = &now
	r.UpdatedAt = now

	return s.repo.Update(id, &r)
}

type AttemptResult struct {
	Success      bool                `json:"success"`
	Rewards      models.RewardsGiven `json:"rewards"`
	RunCompleted bool                `json:"runCompleted"`
}

func (s *Run) AttemptBoss(runID, stepID string, lat, lon float64) (*AttemptResult, error) {
	r, err := s.repo.GetByID(runID)
	if err != nil {
		return nil, errors.New("run not found")
	}

	if r.State != "active" {
		return nil, errors.New("run is not active")
	}

	// Idempotency: check if boss already killed
	for _, ks := range r.KilledSteps {
		if ks.BossStepID == stepID {
			return &AttemptResult{
				Success:      true,
				Rewards:      ks.RewardsGiven,
				RunCompleted: r.State == "completed",
			}, nil
		}
	}

	// Get the boss step
	bs, err := s.bsRepo.GetByID(r.DungeonID, stepID)
	if err != nil {
		return nil, errors.New("boss step not found")
	}

	// Validate step order
	if bs.Order != r.CurrentStep {
		return nil, errors.New("WRONG_STEP_ORDER")
	}

	// Validate location
	distance := functions.HaversineDistance(lat, lon, bs.Location.Lat, bs.Location.Lon)
	if distance > bs.Location.RadiusMeters {
		return nil, errors.New("NOT_IN_RANGE")
	}

	// Roll loot
	rewards := models.RewardsGiven{Gold: bs.GoldReward}
	for _, loot := range bs.LootTable {
		if rand.Float64() <= loot.DropRate {
			qty := loot.MinQty
			if loot.MaxQty > loot.MinQty {
				qty = loot.MinQty + rand.Intn(loot.MaxQty-loot.MinQty+1)
			}
			rewards.Items = append(rewards.Items, models.RewardItem{ItemID: loot.ItemID, Qty: qty})
		}
	}

	// Count total steps to check completion
	totalSteps, err := s.bsRepo.CountByDungeon(r.DungeonID)
	if err != nil {
		return nil, err
	}

	runCompleted := int64(r.CurrentStep) >= totalSteps

	killedStep := models.KilledStep{
		BossStepID:   stepID,
		KilledAt:     time.Now(),
		RewardsGiven: rewards,
	}

	// Execute transactional update in repository
	err = s.repo.ExecuteBossAttempt(context.TODO(), runID, r.PlayerID, rewards, runCompleted, killedStep)
	if err != nil {
		return nil, err
	}

	return &AttemptResult{
		Success:      true,
		Rewards:      rewards,
		RunCompleted: runCompleted,
	}, nil
}

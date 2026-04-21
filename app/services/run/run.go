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
	repo        repositories.RunRepository
	dungeonRepo repositories.DungeonRepository
	bsRepo      repositories.BossStepRepository
	playerRepo  repositories.PlayerRepository
	itemRepo    repositories.ItemRepository
	validate    *validator.Validate
}

func New(
	repo repositories.RunRepository,
	dungeonRepo repositories.DungeonRepository,
	bsRepo repositories.BossStepRepository,
	playerRepo repositories.PlayerRepository,
	itemRepo repositories.ItemRepository,
) *Run {
	return &Run{
		repo:        repo,
		dungeonRepo: dungeonRepo,
		bsRepo:      bsRepo,
		playerRepo:  playerRepo,
		itemRepo:    itemRepo,
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
	r, err := s.repo.GetByID(id)
	if err != nil {
		return r, err
	}
	s.enrichRun(&r)
	return r, nil
}

func (s *Run) enrichRun(r *models.Run) {
	for i := range r.KilledSteps {
		s.enrichRewards(&r.KilledSteps[i].RewardsGiven)
	}
}

func (s *Run) enrichRewards(rg *models.RewardsGiven) {
	for i := range rg.Items {
		item, err := s.itemRepo.GetByID(rg.Items[i].ItemID)
		if err == nil {
			rg.Items[i].ItemDetails = &models.ItemDefResponse{
				ID:          item.CustomID,
				Name:        item.Name,
				Type:        item.Type,
				Rarity:      item.Rarity,
				Description: item.Description,
				Tradable:    item.Tradable,
				BaseValue:   item.BaseValue,
				Stats:       item.Stats,
			}
		}
	}
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

	// Calculate 50% of earned rewards
	var totalGold int64
	itemMap := make(map[string]int)

	for _, ks := range r.KilledSteps {
		totalGold += ks.RewardsGiven.Gold
		for _, it := range ks.RewardsGiven.Items {
			itemMap[it.ItemID] += it.Qty
		}
	}

	// Apply 50% penalty
	finalGold := int64(float64(totalGold) * 0.5)
	var finalItems []models.RewardItem
	for id, qty := range itemMap {
		finalQty := int(float64(qty) * 0.5)
		if finalQty > 0 {
			finalItems = append(finalItems, models.RewardItem{ItemID: id, Qty: finalQty})
		}
	}

	return s.repo.CommitRewards(context.TODO(), id, r.PlayerID, finalGold, finalItems, "abandoned")
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
	radius := float64(500)
	if bs.Location.RadiusMeters != nil {
		radius = *bs.Location.RadiusMeters
	}
	if distance > radius {
		return nil, errors.New("NOT_IN_RANGE")
	}

	// Roll loot for THIS boss
	gold := int64(0)
	if bs.GoldReward != nil {
		gold = *bs.GoldReward
	}
	rewards := models.RewardsGiven{Gold: gold}
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

	// Save progress (without committing to player account yet)
	err = s.repo.ExecuteBossAttempt(context.TODO(), runID, killedStep)
	if err != nil {
		return nil, err
	}

	// If completed, commit ALL rewards + dungeon completion bonus
	if runCompleted {
		d, _ := s.dungeonRepo.GetByID(r.DungeonID)
		
		var totalGold int64 = rewards.Gold // include current boss
		itemMap := make(map[string]int)
		for _, it := range rewards.Items {
			itemMap[it.ItemID] += it.Qty
		}

		// Also aggregate previous steps
		for _, ks := range r.KilledSteps {
			totalGold += ks.RewardsGiven.Gold
			for _, it := range ks.RewardsGiven.Items {
				itemMap[it.ItemID] += it.Qty
			}
		}

		// Add Dungeon Completion Bonus
		totalGold += d.CompletionGoldReward
		for _, loot := range d.CompletionLootTable {
			if rand.Float64() <= loot.DropRate {
				qty := loot.MinQty + rand.Intn(loot.MaxQty-loot.MinQty+1)
				itemMap[loot.ItemID] += qty
			}
		}

		var finalItems []models.RewardItem
		for id, qty := range itemMap {
			finalItems = append(finalItems, models.RewardItem{ItemID: id, Qty: qty})
		}

		err = s.repo.CommitRewards(context.TODO(), runID, r.PlayerID, totalGold, finalItems, "completed")
		if err != nil {
			return nil, err
		}
	}

	s.enrichRewards(&rewards)

	return &AttemptResult{
		Success:      true,
		Rewards:      rewards,
		RunCompleted: runCompleted,
	}, nil
}

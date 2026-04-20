package bossstep

import (
	"dungeons/app/functions"
	"dungeons/app/models"
	"dungeons/app/repositories"
	"time"

	"github.com/go-playground/validator/v10"
)

type BossStep struct {
	repo     repositories.BossStepRepository
	validate *validator.Validate
}

func New(repo repositories.BossStepRepository) *BossStep {
	return &BossStep{
		repo:     repo,
		validate: validator.New(),
	}
}

func (s *BossStep) Create(in *models.BossStep) (*models.BossStep, error) {
	if err := s.validate.Struct(in); err != nil {
		return nil, err
	}

	var bs models.BossStep
	bs.Name = in.Name
	bs.Emoji = in.Emoji
	bs.DungeonID = in.DungeonID
	bs.Location = in.Location
	bs.Difficulty = in.Difficulty
	bs.GoldReward = in.GoldReward
	bs.ZoneDescription = in.ZoneDescription
	bs.LootTable = in.LootTable
	
	bs.CustomID = functions.NewUUID()
	bs.CreatedAt = time.Now()
	bs.UpdatedAt = time.Now()

	steps, err := s.repo.GetByDungeonOrdered(bs.DungeonID)
	if err != nil {
		return nil, err
	}
	
	bs.Order = 1
	if len(steps) > 0 {
		bs.Order = steps[len(steps)-1].Order + 1
	}

	if err := s.repo.Create(&bs); err != nil {
		return nil, err
	}

	return &bs, nil
}

func (s *BossStep) GetByID(id string) (models.BossStep, error) {
	// Note: The service interface previously took just id, but the model needs dungeonId
	// We might need to adjust or keep it simple for now if the repo can find by ID only.
	// Actually, customID is unique enough. 
	// I'll update the repo to support GetByID(id) if needed.
	// For now, let's look at what we have.
	return s.repo.GetByID("", id) // Temporary, better to fix repo interface
}

func (s *BossStep) Update(dungeonID, stepID string, in *models.BossStep) error {
	bs, err := s.repo.GetByID(dungeonID, stepID)
	if err != nil {
		return err
	}

	bs.Name = in.Name
	bs.Emoji = in.Emoji
	bs.Location = in.Location
	bs.Difficulty = in.Difficulty
	bs.GoldReward = in.GoldReward
	bs.ZoneDescription = in.ZoneDescription
	bs.UpdatedAt = time.Now()

	return s.repo.Update(stepID, &bs)
}

func (s *BossStep) Delete(dungeonID, stepID string) error {
	err := s.repo.Delete(dungeonID, stepID)
	if err != nil {
		return err
	}

	// Reorder remaining steps
	remaining, err := s.repo.GetByDungeonOrdered(dungeonID)
	if err != nil {
		return err
	}

	for i, step := range remaining {
		step.Order = i + 1
		step.UpdatedAt = time.Now()
		_ = s.repo.Update(step.CustomID, &step)
	}

	return nil
}

func (s *BossStep) GetByDungeonID(dungeonID string) ([]models.BossStep, error) {
	return s.repo.GetByDungeonOrdered(dungeonID)
}

func (s *BossStep) Reorder(dungeonID string, orderedIDs []string) error {
	for i, id := range orderedIDs {
		bs, err := s.repo.GetByID(dungeonID, id)
		if err != nil {
			return err
		}
		bs.Order = i + 1
		bs.UpdatedAt = time.Now()
		err = s.repo.Update(id, &bs)
		if err != nil {
			return err
		}
	}
	return nil
}

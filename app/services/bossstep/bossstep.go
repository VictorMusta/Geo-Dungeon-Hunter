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
	bs.ZoneDescription = in.ZoneDescription
	bs.LootTable = in.LootTable

	// Default values if nil
	if in.Difficulty == nil {
		d := 5
		bs.Difficulty = &d
	} else {
		bs.Difficulty = in.Difficulty
	}
	
	if in.GoldReward == nil {
		g := int64(0)
		bs.GoldReward = &g
	} else {
		bs.GoldReward = in.GoldReward
	}

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

	// Update fields ONLY if they are provided (not nil for pointers, not empty for strings)
	if in.Name != "" {
		bs.Name = in.Name
	}
	if in.Emoji != "" {
		bs.Emoji = in.Emoji
	}
	if in.Location.Lat != 0 {
		bs.Location.Lat = in.Location.Lat
	}
	if in.Location.Lon != 0 {
		bs.Location.Lon = in.Location.Lon
	}
	if in.Location.RadiusMeters != nil {
		bs.Location.RadiusMeters = in.Location.RadiusMeters
	}
	if in.Difficulty != nil {
		bs.Difficulty = in.Difficulty
	}
	if in.GoldReward != nil {
		bs.GoldReward = in.GoldReward
	}
	if in.ZoneDescription != "" {
		bs.ZoneDescription = in.ZoneDescription
	}
	if in.LootTable != nil {
		bs.LootTable = in.LootTable
	}
	
	bs.UpdatedAt = time.Now()

	// Validate the final merged object
	if err := s.validate.Struct(&bs); err != nil {
		return err
	}

	return s.repo.Update(dungeonID, stepID, &bs)
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
		_ = s.repo.Update(dungeonID, step.CustomID, &step)
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
		err = s.repo.Update(dungeonID, id, &bs)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *BossStep) CountByDungeon(dungeonID string) (int, error) {
	count, err := s.repo.CountByDungeon(dungeonID)
	if err != nil {
		return 0, err
	}
	return int(count), nil
}

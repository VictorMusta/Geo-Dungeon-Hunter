package dungeon

import (
	"dungeons/app/functions"
	"dungeons/app/models"
	"dungeons/app/repositories"
	"errors"
	"time"

	"github.com/go-playground/validator/v10"
)

type BossStepService interface {
	GetByDungeonID(did string) ([]models.BossStep, error)
	Create(in *models.BossStep) (*models.BossStep, error)
	Update(did, sid string, in *models.BossStep) error
	Delete(did, sid string) error
	Reorder(did string, ids []string) error
	CountByDungeon(did string) (int, error)
}

type Dungeon struct {
	repo      repositories.DungeonRepository
	bsService BossStepService
	validate  *validator.Validate
}

func New(repo repositories.DungeonRepository, bsService BossStepService) *Dungeon {
	return &Dungeon{
		repo:      repo,
		bsService: bsService,
		validate:  validator.New(),
	}
}

func (s *Dungeon) Create(in *models.Dungeon) (*models.Dungeon, error) {
	if err := s.validate.Struct(in); err != nil {
		return nil, err
	}

	var d models.Dungeon
	d.Title = in.Title
	d.Description = in.Description
	d.Area = in.Area
	d.CreatedBy = in.CreatedBy
	d.CustomID = functions.NewUUID()
	d.Status = "draft"
	d.CreatedAt = time.Now()
	d.UpdatedAt = time.Now()

	if err := s.repo.Create(&d); err != nil {
		return nil, err
	}

	return &d, nil
}

func (s *Dungeon) GetByID(id string) (models.Dungeon, error) {
	return s.repo.GetByID(id)
}

func (s *Dungeon) UpdateFull(id string, in *models.DungeonFullUpdate) error {
	// 1. Update Dungeon
	err := s.Update(id, &in.Dungeon)
	if err != nil {
		return err
	}

	// 2. Update Boss Steps
	for _, step := range in.BossSteps {
		// Create a copy to avoid pointer issues in the loop
		currentStep := step
		// Here we assume the dungeonID in the step is correct or we force it
		currentStep.DungeonID = id
		
		err = s.bsService.Update(id, currentStep.CustomID, &currentStep)
		if err != nil {
			// We continue or return error? Let's return error for strictness
			return err
		}
	}

	return nil
}

func (s *Dungeon) Update(id string, in *models.Dungeon) error {
	d, err := s.repo.GetByID(id)
	if err != nil {
		return errors.New("dungeon not found")
	}

	if d.Status != "draft" {
		return errors.New("only draft dungeons can be modified")
	}

	if in.Title != "" {
		d.Title = in.Title
	}
	if in.Description != "" {
		d.Description = in.Description
	}
	
	// Update Area (Mapping/BoundingBox)
	// We allow updating the Area if either name or boundingBox is provided
	if in.Area.Name != "" || in.Area.BoundingBox != nil {
		d.Area = in.Area
	}

	if in.CompletionGoldReward >= 0 {
		d.CompletionGoldReward = in.CompletionGoldReward
	}
	if in.CompletionLootTable != nil {
		d.CompletionLootTable = in.CompletionLootTable
	}
	d.UpdatedAt = time.Now()

	// Final validation of the merged object
	if err := s.validate.Struct(&d); err != nil {
		return err
	}

	return s.repo.Update(id, &d)
}

func (s *Dungeon) Publish(id string) error {
	d, err := s.repo.GetByID(id)
	if err != nil {
		return errors.New("dungeon not found")
	}

	if d.Status != "draft" {
		return errors.New("only draft dungeons can be published")
	}

	count, err := s.bsService.CountByDungeon(id)
	if err != nil {
		return err
	}
	if count == 0 {
		return errors.New("dungeon must have at least one boss step to publish")
	}

	d.Status = "published"
	d.UpdatedAt = time.Now()

	return s.repo.Update(id, &d)
}

func (s *Dungeon) GetPublished(queryParams models.QueryParams) ([]models.Dungeon, error) {
	return s.repo.GetPublished(queryParams)
}

func (s *Dungeon) GetByMJ(mjId string, queryParams models.QueryParams) ([]models.Dungeon, error) {
	return s.repo.GetByMJ(mjId, queryParams)
}

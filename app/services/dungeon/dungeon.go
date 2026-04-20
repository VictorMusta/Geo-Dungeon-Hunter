package dungeon

import (
	"dungeons/app/functions"
	"dungeons/app/models"
	"dungeons/app/repositories"
	"errors"
	"time"

	"github.com/go-playground/validator/v10"
)

type Dungeon struct {
	repo     repositories.DungeonRepository
	bsRepo   repositories.BossStepRepository
	validate *validator.Validate
}

func New(repo repositories.DungeonRepository, bsRepo repositories.BossStepRepository) *Dungeon {
	return &Dungeon{
		repo:     repo,
		bsRepo:   bsRepo,
		validate: validator.New(),
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
	if in.Area.Name != "" {
		d.Area = in.Area
	}
	d.UpdatedAt = time.Now()

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

	count, err := s.bsRepo.CountByDungeon(id)
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

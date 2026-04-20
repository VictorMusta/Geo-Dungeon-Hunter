package item

import (
	"dungeons/app/functions"
	"dungeons/app/models"
	"dungeons/app/repositories"
	"time"

	"github.com/go-playground/validator/v10"
)

type Item struct {
	repo     repositories.ItemRepository
	validate *validator.Validate
}

func New(repo repositories.ItemRepository) *Item {
	return &Item{
		repo:     repo,
		validate: validator.New(),
	}
}

func (s *Item) Get(queryParams models.QueryParams) ([]models.ItemDef, error) {
	return s.repo.Get(queryParams)
}

func (s *Item) Create(in *models.ItemDef) (*models.ItemDef, error) {
	if err := s.validate.Struct(in); err != nil {
		return nil, err
	}

	var item models.ItemDef
	item.Name = in.Name
	item.Type = in.Type
	item.Rarity = in.Rarity
	item.Description = in.Description
	item.Stats = in.Stats
	item.Tradable = in.Tradable
	item.BaseValue = in.BaseValue

	item.CustomID = functions.NewUUID()
	item.CreatedAt = time.Now()
	item.UpdatedAt = time.Now()

	if err := s.repo.Create(&item); err != nil {
		return nil, err
	}

	return &item, nil
}

func (s *Item) GetByID(id string) (models.ItemDef, error) {
	return s.repo.GetByID(id)
}

func (s *Item) Update(id string, in *models.ItemDef) error {
	if err := s.validate.Struct(in); err != nil {
		return err
	}

	item, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}

	item.Name = in.Name
	item.Type = in.Type
	item.Rarity = in.Rarity
	item.Description = in.Description
	item.Tradable = in.Tradable
	item.BaseValue = in.BaseValue
	item.UpdatedAt = time.Now()

	return s.repo.Update(id, &item)
}

func (s *Item) Delete(id string) error {
	return s.repo.Delete(id)
}

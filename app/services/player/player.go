package player

import (
	"dungeons/app/functions"
	"dungeons/app/models"
	"dungeons/app/repositories"
	"errors"
	"time"

	"github.com/go-playground/validator/v10"
)

type Player struct {
	repo     repositories.PlayerRepository
	validate *validator.Validate
}

func New(repo repositories.PlayerRepository) *Player {
	return &Player{
		repo:     repo,
		validate: validator.New(),
	}
}

// Get services to get list of Player
func (p *Player) Get(queryParams models.QueryParams) ([]models.Player, error) {
	return p.repo.Get(queryParams)
}

// Create new player
func (p *Player) Create(in *models.Player) (*models.Player, error) {
	// Check input fields
	err := p.validate.Struct(in)
	if err != nil {
		return nil, err
	}

	// Check if pseudo already exists (Get-or-Create behavior for UX)
	existing, err := p.repo.FindByDisplayName(in.DisplayName)
	if err == nil {
		return &existing, nil
	}

	var player models.Player
	player.DisplayName = in.DisplayName
	player.CustomID = functions.NewUUID()
	player.CreatedAt = time.Now()

	// Hash the password if provided
	if in.Password != "" {
		hashed, err := functions.HashAndSalt(in.Password)
		if err != nil {
			return nil, err
		}
		player.Password = string(hashed)
	}

	err = p.repo.Create(&player)
	if err != nil {
		return nil, err
	}

	return &player, nil
}

// GetByID to get one Player by ID
func (p *Player) GetByID(id string) (models.Player, error) {
	return p.repo.GetByID(id)
}

// Update to update a Player
func (p *Player) Update(id string, in *models.Player) error {
	// Check input fields
	err := p.validate.Struct(in)
	if err != nil {
		return err
	}

	player, err := p.repo.GetByID(id)
	if err != nil {
		return err
	}

	if in.DisplayName != "" {
		player.DisplayName = in.DisplayName
	}
	player.UpdatedAt = time.Now()

	return p.repo.Update(id, &player)
}

// Suspend to suspend a Player
func (p *Player) Suspend(id string) error {
	return p.repo.Suspend(id)
}

// GetByIds to get list of Player by Ids
func (p *Player) GetByIds(ids []string) ([]models.Player, error) {
	return p.repo.GetByIDs(ids)
}

// Login verifies credentials and returns a token
func (p *Player) Login(displayName, password string) (*models.LoginResponse, error) {
	player, err := p.repo.FindByDisplayName(displayName)
	if err != nil {
		return nil, errors.New("identifiant ou mot de passe incorrect")
	}

	// Double check password
	err = functions.CheckPassword(password, player.Password)
	if err != nil {
		return nil, errors.New("identifiant ou mot de passe incorrect")
	}

	// Generate JWT
	token, err := functions.GenerateToken(player.CustomID)
	if err != nil {
		return nil, err
	}

	return &models.LoginResponse{
		Token:  token,
		Player: player,
	}, nil
}

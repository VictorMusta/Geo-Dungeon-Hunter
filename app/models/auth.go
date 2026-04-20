package models

type LoginRequest struct {
	DisplayName string `json:"display_name" validate:"required"`
	Password    string `json:"password" validate:"required"`
}

type LoginResponse struct {
	Token  string `json:"token"`
	Player Player `json:"player"`
}

package models

import (
	"time"
)

type User struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	Email         string    `gorm:"unique;not null" json:"email"`
	Password      string    `gorm:"not null" json:"-"`
	Name          string    `gorm:"not null" json:"name"`
	BirthDate     time.Time `json:"birth_date"`
	HeightCm      float64   `json:"height_cm"`
	WeightKg      float64   `json:"weight_kg"`
	ActivityLevel string    `gorm:"default:'sedentary'" json:"activity_level"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	Name     string `json:"name" binding:"required"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}

type UpdateProfileRequest struct {
	Name          string    `json:"name"`
	BirthDate     time.Time `json:"birth_date"`
	HeightCm      float64   `json:"height_cm"`
	WeightKg      float64   `json:"weight_kg"`
	ActivityLevel string    `json:"activity_level"`
}

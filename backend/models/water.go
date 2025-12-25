package models

import "time"

// WaterIntake represents daily water intake tracking
type WaterIntake struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	UserID    uint      `json:"user_id" gorm:"not null"`
	Glasses   int       `json:"glasses" gorm:"default:0"` // Number of glasses (1 glass = 250ml)
	Goal      int       `json:"goal" gorm:"default:8"`    // Daily goal in glasses
	Date      string    `json:"date" gorm:"size:10;not null"` // Format: YYYY-MM-DD
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// WaterIntakeResponse is the response structure for water intake
type WaterIntakeResponse struct {
	ID         uint    `json:"id"`
	Glasses    int     `json:"glasses"`
	Goal       int     `json:"goal"`
	Date       string  `json:"date"`
	Percentage float64 `json:"percentage"`
	Remaining  int     `json:"remaining"`
}

// GetPercentage calculates the percentage of goal achieved
func (w *WaterIntake) GetPercentage() float64 {
	if w.Goal == 0 {
		return 0
	}
	percentage := float64(w.Glasses) / float64(w.Goal) * 100
	if percentage > 100 {
		return 100
	}
	return percentage
}

// GetRemaining calculates remaining glasses to reach goal
func (w *WaterIntake) GetRemaining() int {
	remaining := w.Goal - w.Glasses
	if remaining < 0 {
		return 0
	}
	return remaining
}

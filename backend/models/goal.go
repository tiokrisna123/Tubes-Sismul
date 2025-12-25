package models

import "time"

// Goal represents a health goal
type Goal struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	UserID      uint      `json:"user_id" gorm:"not null"`
	Title       string    `json:"title" gorm:"size:200;not null"`
	Description string    `json:"description" gorm:"size:500"`
	Type        string    `json:"type" gorm:"size:50;not null"` // weight, exercise, water, sleep, custom
	Target      float64   `json:"target" gorm:"not null"`
	Current     float64   `json:"current" gorm:"default:0"`
	Unit        string    `json:"unit" gorm:"size:20"` // kg, minutes, glasses, hours, etc.
	Deadline    string    `json:"deadline" gorm:"size:10"` // Format: YYYY-MM-DD
	IsCompleted bool      `json:"is_completed" gorm:"default:false"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// GoalResponse is the response structure for a goal
type GoalResponse struct {
	ID          uint    `json:"id"`
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Type        string  `json:"type"`
	Target      float64 `json:"target"`
	Current     float64 `json:"current"`
	Unit        string  `json:"unit"`
	Deadline    string  `json:"deadline"`
	IsCompleted bool    `json:"is_completed"`
	Progress    float64 `json:"progress"` // percentage
	DaysLeft    int     `json:"days_left"`
}

// GetProgress calculates the progress percentage
func (g *Goal) GetProgress() float64 {
	if g.Target == 0 {
		return 0
	}
	progress := g.Current / g.Target * 100
	if progress > 100 {
		return 100
	}
	return progress
}

// GoalType constants
const (
	GoalTypeWeight   = "weight"
	GoalTypeExercise = "exercise"
	GoalTypeWater    = "water"
	GoalTypeSleep    = "sleep"
	GoalTypeCustom   = "custom"
)

// GetGoalTypeIcon returns icon for goal type
func GetGoalTypeIcon(goalType string) string {
	icons := map[string]string{
		GoalTypeWeight:   "⚖️",
		GoalTypeExercise: "🏃",
		GoalTypeWater:    "💧",
		GoalTypeSleep:    "😴",
		GoalTypeCustom:   "🎯",
	}
	if icon, ok := icons[goalType]; ok {
		return icon
	}
	return "🎯"
}

package models

import "time"

// Reminder represents a health reminder
type Reminder struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	UserID    uint      `json:"user_id" gorm:"not null"`
	Type      string    `json:"type" gorm:"size:50;not null"` // water, meal, exercise, meditation, rest, custom
	Label     string    `json:"label" gorm:"size:200;not null"`
	Time      string    `json:"time" gorm:"size:10;not null"` // Format: HH:MM
	IsActive  bool      `json:"is_active" gorm:"default:true"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ReminderType constants
const (
	ReminderTypeWater      = "water"
	ReminderTypeMeal       = "meal"
	ReminderTypeExercise   = "exercise"
	ReminderTypeMeditation = "meditation"
	ReminderTypeRest       = "rest"
	ReminderTypeCustom     = "custom"
)

// CreateReminderRequest is the request structure for creating a reminder
type CreateReminderRequest struct {
	Type  string `json:"type" binding:"required"`
	Label string `json:"label" binding:"required"`
	Time  string `json:"time" binding:"required"`
}

// UpdateReminderRequest is the request structure for updating a reminder
type UpdateReminderRequest struct {
	Type     string `json:"type"`
	Label    string `json:"label"`
	Time     string `json:"time"`
	IsActive *bool  `json:"is_active"`
}

// ReminderResponse is the response structure for a reminder
type ReminderResponse struct {
	ID       uint   `json:"id"`
	Type     string `json:"type"`
	Label    string `json:"label"`
	Time     string `json:"time"`
	IsActive bool   `json:"is_active"`
	Icon     string `json:"icon"`
}

// GetReminderIcon returns icon for reminder type
func GetReminderIcon(reminderType string) string {
	icons := map[string]string{
		ReminderTypeWater:      "💧",
		ReminderTypeMeal:       "🍽️",
		ReminderTypeExercise:   "🏃",
		ReminderTypeMeditation: "🧘",
		ReminderTypeRest:       "😴",
		ReminderTypeCustom:     "⏰",
	}
	if icon, ok := icons[reminderType]; ok {
		return icon
	}
	return "⏰"
}

// ToResponse converts Reminder to ReminderResponse
func (r *Reminder) ToResponse() ReminderResponse {
	return ReminderResponse{
		ID:       r.ID,
		Type:     r.Type,
		Label:    r.Label,
		Time:     r.Time,
		IsActive: r.IsActive,
		Icon:     GetReminderIcon(r.Type),
	}
}

// DefaultReminders returns default reminders for new users
func DefaultReminders() []Reminder {
	return []Reminder{
		{Type: ReminderTypeWater, Label: "Minum Air", Time: "08:00", IsActive: true},
		{Type: ReminderTypeMeal, Label: "Sarapan Sehat", Time: "08:30", IsActive: true},
		{Type: ReminderTypeExercise, Label: "Olahraga Pagi", Time: "07:00", IsActive: true},
		{Type: ReminderTypeMeditation, Label: "Meditasi", Time: "06:30", IsActive: true},
		{Type: ReminderTypeWater, Label: "Minum Air", Time: "12:00", IsActive: true},
		{Type: ReminderTypeMeal, Label: "Makan Siang", Time: "12:30", IsActive: true},
		{Type: ReminderTypeRest, Label: "Istirahat Siang", Time: "14:00", IsActive: true},
		{Type: ReminderTypeWater, Label: "Minum Air", Time: "16:00", IsActive: true},
		{Type: ReminderTypeMeal, Label: "Makan Malam", Time: "19:00", IsActive: true},
		{Type: ReminderTypeRest, Label: "Persiapan Tidur", Time: "21:00", IsActive: true},
	}
}

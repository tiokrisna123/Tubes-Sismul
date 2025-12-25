package models

import (
	"time"
)

type HealthData struct {
	ID             uint      `gorm:"primaryKey" json:"id"`
	UserID         uint      `gorm:"not null" json:"user_id"`
	WeightKg       float64   `json:"weight_kg"`
	HeightCm       float64   `json:"height_cm"`
	BMI            float64   `json:"bmi"`
	ActivityLevel  string    `json:"activity_level"`
	EmotionalState string    `json:"emotional_state"`
	DailySchedule  string    `gorm:"type:text" json:"daily_schedule"`
	Notes          string    `json:"notes"`
	RecordDate     time.Time `json:"record_date"`
	CreatedAt      time.Time `json:"created_at"`
}

type HealthDataRequest struct {
	WeightKg       float64 `json:"weight_kg" binding:"required"`
	HeightCm       float64 `json:"height_cm" binding:"required"`
	ActivityLevel  string  `json:"activity_level"`
	EmotionalState string  `json:"emotional_state"`
	DailySchedule  string  `json:"daily_schedule"`
	Notes          string  `json:"notes"`
}

type DashboardData struct {
	LatestHealth    *HealthData          `json:"latest_health"`
	BMICategory     string               `json:"bmi_category"`
	HealthScore     int                  `json:"health_score"`
	TotalRecords    int64                `json:"total_records"`
	RecentSymptoms  []Symptom            `json:"recent_symptoms"`
	WeeklyProgress  []HealthData         `json:"weekly_progress"`
	Recommendations []RecommendationItem `json:"recommendations"`
}

type RecommendationItem struct {
	Type        string `json:"type"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Priority    string `json:"priority"`
}

func CalculateBMI(weightKg, heightCm float64) float64 {
	if heightCm <= 0 {
		return 0
	}
	heightM := heightCm / 100
	return weightKg / (heightM * heightM)
}

func GetBMICategory(bmi float64) string {
	switch {
	case bmi < 18.5:
		return "Underweight"
	case bmi < 25:
		return "Normal"
	case bmi < 30:
		return "Overweight"
	default:
		return "Obese"
	}
}

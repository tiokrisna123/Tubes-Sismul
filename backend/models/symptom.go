package models

import (
	"time"
)

type Symptom struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	UserID      uint      `gorm:"not null" json:"user_id"`
	SymptomType string    `gorm:"not null" json:"symptom_type"` // physical, mental
	SymptomName string    `gorm:"not null" json:"symptom_name"`
	Severity    int       `json:"severity"` // 1-10
	Notes       string    `json:"notes"`
	LoggedAt    time.Time `json:"logged_at"`
}

type SymptomRequest struct {
	SymptomType string `json:"symptom_type" binding:"required"`
	SymptomName string `json:"symptom_name" binding:"required"`
	Severity    int    `json:"severity" binding:"required,min=1,max=10"`
	Notes       string `json:"notes"`
}

type SymptomTemplate struct {
	ID          uint   `gorm:"primaryKey" json:"id"`
	SymptomType string `json:"symptom_type"`
	SymptomName string `json:"symptom_name"`
	Description string `json:"description"`
}

// Predefined symptoms
var PhysicalSymptoms = []string{
	"Demam",
	"Flu",
	"Batuk",
	"Pilek",
	"Sakit Kepala",
	"Tekanan Darah Tinggi",
	"Kolesterol Tinggi",
	"Maag",
	"Gangguan Pencernaan",
	"Nyeri Otot",
	"Kelelahan Fisik",
	"Obesitas",
	"Nyeri Sendi",
	"Sesak Napas",
	"Pusing",
}

var MentalSymptoms = []string{
	"Stres",
	"Kecemasan",
	"Depresi Ringan",
	"Mudah Marah",
	"Gangguan Tidur",
	"Burnout",
	"Kesepian Sosial",
	"Sulit Konsentrasi",
	"Mood Swing",
	"Overthinking",
}

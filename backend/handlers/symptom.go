package handlers

import (
	"net/http"
	"time"

	"health-tracker/database"
	"health-tracker/models"
	"health-tracker/utils"

	"github.com/gin-gonic/gin"
)

// GetSymptomList returns all available symptom templates
func GetSymptomList(c *gin.Context) {
	var symptoms []models.SymptomTemplate
	database.DB.Find(&symptoms)

	// Group by type
	physical := []models.SymptomTemplate{}
	mental := []models.SymptomTemplate{}

	for _, s := range symptoms {
		if s.SymptomType == "physical" {
			physical = append(physical, s)
		} else {
			mental = append(mental, s)
		}
	}

	utils.SuccessResponse(c, http.StatusOK, "Symptom list retrieved", gin.H{
		"physical": physical,
		"mental":   mental,
	})
}

// LogSymptom records a new symptom
func LogSymptom(c *gin.Context) {
	userID := c.GetUint("userID")

	var req models.SymptomRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request: "+err.Error())
		return
	}

	symptom := models.Symptom{
		UserID:      userID,
		SymptomType: req.SymptomType,
		SymptomName: req.SymptomName,
		Severity:    req.Severity,
		Notes:       req.Notes,
		LoggedAt:    time.Now(),
	}

	if result := database.DB.Create(&symptom); result.Error != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to log symptom")
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "Symptom logged successfully", symptom)
}

// LogMultipleSymptoms records multiple symptoms at once
func LogMultipleSymptoms(c *gin.Context) {
	userID := c.GetUint("userID")

	var requests []models.SymptomRequest
	if err := c.ShouldBindJSON(&requests); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request: "+err.Error())
		return
	}

	symptoms := make([]models.Symptom, len(requests))
	now := time.Now()

	for i, req := range requests {
		symptoms[i] = models.Symptom{
			UserID:      userID,
			SymptomType: req.SymptomType,
			SymptomName: req.SymptomName,
			Severity:    req.Severity,
			Notes:       req.Notes,
			LoggedAt:    now,
		}
	}

	if result := database.DB.Create(&symptoms); result.Error != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to log symptoms")
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "Symptoms logged successfully", symptoms)
}

// GetSymptomHistory returns user's symptom history
func GetSymptomHistory(c *gin.Context) {
	userID := c.GetUint("userID")

	var symptoms []models.Symptom
	database.DB.Where("user_id = ?", userID).Order("logged_at desc").Limit(50).Find(&symptoms)

	// Group by date
	grouped := make(map[string][]models.Symptom)
	for _, s := range symptoms {
		date := s.LoggedAt.Format("2006-01-02")
		grouped[date] = append(grouped[date], s)
	}

	utils.SuccessResponse(c, http.StatusOK, "Symptom history retrieved", gin.H{
		"symptoms": symptoms,
		"grouped":  grouped,
	})
}

// GetSymptomStats returns symptom statistics
func GetSymptomStats(c *gin.Context) {
	userID := c.GetUint("userID")

	// Most frequent symptoms
	type SymptomCount struct {
		SymptomName string `json:"symptom_name"`
		Count       int    `json:"count"`
	}

	var frequentSymptoms []SymptomCount
	database.DB.Model(&models.Symptom{}).
		Select("symptom_name, COUNT(*) as count").
		Where("user_id = ?", userID).
		Group("symptom_name").
		Order("count desc").
		Limit(5).
		Scan(&frequentSymptoms)

	// Symptoms this week
	weekAgo := time.Now().AddDate(0, 0, -7)
	var weekCount int64
	database.DB.Model(&models.Symptom{}).
		Where("user_id = ? AND logged_at > ?", userID, weekAgo).
		Count(&weekCount)

	// Average severity
	var avgSeverity float64
	database.DB.Model(&models.Symptom{}).
		Select("COALESCE(AVG(severity), 0)").
		Where("user_id = ?", userID).
		Scan(&avgSeverity)

	utils.SuccessResponse(c, http.StatusOK, "Symptom stats retrieved", gin.H{
		"frequent_symptoms":  frequentSymptoms,
		"symptoms_this_week": weekCount,
		"average_severity":   avgSeverity,
	})
}

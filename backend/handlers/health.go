package handlers

import (
	"net/http"
	"time"

	"health-tracker/database"
	"health-tracker/models"
	"health-tracker/utils"

	"github.com/gin-gonic/gin"
)

// CreateHealthData submits new health data
func CreateHealthData(c *gin.Context) {
	userID := c.GetUint("userID")

	var req models.HealthDataRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request: "+err.Error())
		return
	}

	// Calculate BMI
	bmi := models.CalculateBMI(req.WeightKg, req.HeightCm)

	healthData := models.HealthData{
		UserID:         userID,
		WeightKg:       req.WeightKg,
		HeightCm:       req.HeightCm,
		BMI:            bmi,
		ActivityLevel:  req.ActivityLevel,
		EmotionalState: req.EmotionalState,
		DailySchedule:  req.DailySchedule,
		Notes:          req.Notes,
		RecordDate:     time.Now(),
	}

	if result := database.DB.Create(&healthData); result.Error != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to save health data")
		return
	}

	// Also update user's base info
	database.DB.Model(&models.User{}).Where("id = ?", userID).Updates(map[string]interface{}{
		"weight_kg":      req.WeightKg,
		"height_cm":      req.HeightCm,
		"activity_level": req.ActivityLevel,
	})

	utils.SuccessResponse(c, http.StatusCreated, "Health data saved", gin.H{
		"health_data":  healthData,
		"bmi_category": models.GetBMICategory(bmi),
	})
}

// GetHealthData returns all health records for current user
func GetHealthData(c *gin.Context) {
	userID := c.GetUint("userID")

	var healthData []models.HealthData
	database.DB.Where("user_id = ?", userID).Order("record_date desc").Find(&healthData)

	utils.SuccessResponse(c, http.StatusOK, "Health data retrieved", healthData)
}

// GetLatestHealthData returns the latest health record
func GetLatestHealthData(c *gin.Context) {
	userID := c.GetUint("userID")

	var healthData models.HealthData
	result := database.DB.Where("user_id = ?", userID).Order("record_date desc").First(&healthData)

	if result.Error != nil {
		utils.SuccessResponse(c, http.StatusOK, "No health data found", nil)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Latest health data", gin.H{
		"health_data":  healthData,
		"bmi_category": models.GetBMICategory(healthData.BMI),
	})
}

// GetDashboard returns dashboard summary data
func GetDashboard(c *gin.Context) {
	userID := c.GetUint("userID")

	// Get latest health data
	var latestHealth models.HealthData
	database.DB.Where("user_id = ?", userID).Order("record_date desc").First(&latestHealth)

	// Get total records count
	var totalRecords int64
	database.DB.Model(&models.HealthData{}).Where("user_id = ?", userID).Count(&totalRecords)

	// Get recent symptoms (last 7 days)
	var recentSymptoms []models.Symptom
	weekAgo := time.Now().AddDate(0, 0, -7)
	database.DB.Where("user_id = ? AND logged_at > ?", userID, weekAgo).Order("logged_at desc").Find(&recentSymptoms)

	// Get weekly progress (last 7 records)
	var weeklyProgress []models.HealthData
	database.DB.Where("user_id = ?", userID).Order("record_date desc").Limit(7).Find(&weeklyProgress)

	// Calculate health score (simplified)
	healthScore := calculateHealthScore(latestHealth, recentSymptoms)

	// Get recommendations
	recommendations := getQuickRecommendations(latestHealth, recentSymptoms)

	dashboard := models.DashboardData{
		LatestHealth:    &latestHealth,
		BMICategory:     models.GetBMICategory(latestHealth.BMI),
		HealthScore:     healthScore,
		TotalRecords:    totalRecords,
		RecentSymptoms:  recentSymptoms,
		WeeklyProgress:  weeklyProgress,
		Recommendations: recommendations,
	}

	utils.SuccessResponse(c, http.StatusOK, "Dashboard data retrieved", dashboard)
}

// GetHealthGraph returns graph data for specified period
func GetHealthGraph(c *gin.Context) {
	userID := c.GetUint("userID")
	period := c.Param("period") // week, month, year

	var days int
	switch period {
	case "week":
		days = 7
	case "month":
		days = 30
	case "year":
		days = 365
	default:
		days = 7
	}

	startDate := time.Now().AddDate(0, 0, -days)

	var healthData []models.HealthData
	database.DB.Where("user_id = ? AND record_date >= ?", userID, startDate).
		Order("record_date asc").Find(&healthData)

	// Prepare graph data
	graphData := make([]map[string]interface{}, len(healthData))
	for i, hd := range healthData {
		graphData[i] = map[string]interface{}{
			"date":            hd.RecordDate.Format("2006-01-02"),
			"weight":          hd.WeightKg,
			"bmi":             hd.BMI,
			"emotional_state": hd.EmotionalState,
		}
	}

	utils.SuccessResponse(c, http.StatusOK, "Graph data retrieved", graphData)
}

func calculateHealthScore(health models.HealthData, symptoms []models.Symptom) int {
	score := 100

	// Deduct points based on BMI
	bmiCategory := models.GetBMICategory(health.BMI)
	switch bmiCategory {
	case "Underweight":
		score -= 15
	case "Overweight":
		score -= 10
	case "Obese":
		score -= 25
	}

	// Deduct points for symptoms
	score -= len(symptoms) * 5

	// Deduct based on emotional state
	switch health.EmotionalState {
	case "stressed", "anxious":
		score -= 10
	case "sad":
		score -= 15
	}

	if score < 0 {
		score = 0
	}

	return score
}

func getQuickRecommendations(health models.HealthData, symptoms []models.Symptom) []models.RecommendationItem {
	var recommendations []models.RecommendationItem

	// BMI-based recommendation
	bmiCategory := models.GetBMICategory(health.BMI)
	if bmiCategory != "Normal" {
		recommendations = append(recommendations, models.RecommendationItem{
			Type:        "health",
			Title:       "Perhatikan BMI Anda",
			Description: "BMI Anda termasuk " + bmiCategory + ". Pertimbangkan untuk menyesuaikan pola makan dan olahraga.",
			Priority:    "high",
		})
	}

	// Emotional state recommendation
	if health.EmotionalState == "stressed" || health.EmotionalState == "anxious" {
		recommendations = append(recommendations, models.RecommendationItem{
			Type:        "emotional",
			Title:       "Kelola Stres Anda",
			Description: "Coba teknik relaksasi seperti meditasi atau pernapasan dalam.",
			Priority:    "medium",
		})
	}

	// Symptom-based
	if len(symptoms) > 3 {
		recommendations = append(recommendations, models.RecommendationItem{
			Type:        "health",
			Title:       "Banyak Gejala Terdeteksi",
			Description: "Anda memiliki beberapa gejala. Pertimbangkan untuk berkonsultasi dengan dokter.",
			Priority:    "high",
		})
	}

	return recommendations
}

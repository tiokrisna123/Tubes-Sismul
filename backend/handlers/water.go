package handlers

import (
	"health-tracker/database"
	"health-tracker/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// GetWaterIntake returns today's water intake for the user
func GetWaterIntake(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	today := time.Now().Format("2006-01-02")
	var water models.WaterIntake

	if err := database.DB.Where("user_id = ? AND date = ?", userID, today).First(&water).Error; err != nil {
		// Create new record for today
		water = models.WaterIntake{
			UserID:    userID.(uint),
			Glasses:   0,
			Goal:      8,
			Date:      today,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		database.DB.Create(&water)
	}

	c.JSON(http.StatusOK, models.WaterIntakeResponse{
		ID:         water.ID,
		Glasses:    water.Glasses,
		Goal:       water.Goal,
		Date:       water.Date,
		Percentage: water.GetPercentage(),
		Remaining:  water.GetRemaining(),
	})
}

// AddWaterGlass adds a glass of water
func AddWaterGlass(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	today := time.Now().Format("2006-01-02")
	var water models.WaterIntake

	if err := database.DB.Where("user_id = ? AND date = ?", userID, today).First(&water).Error; err != nil {
		// Create new record for today
		water = models.WaterIntake{
			UserID:    userID.(uint),
			Glasses:   1,
			Goal:      8,
			Date:      today,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		database.DB.Create(&water)
	} else {
		// Update existing record
		water.Glasses++
		water.UpdatedAt = time.Now()
		database.DB.Save(&water)
	}

	c.JSON(http.StatusOK, models.WaterIntakeResponse{
		ID:         water.ID,
		Glasses:    water.Glasses,
		Goal:       water.Goal,
		Date:       water.Date,
		Percentage: water.GetPercentage(),
		Remaining:  water.GetRemaining(),
	})
}

// RemoveWaterGlass removes a glass of water
func RemoveWaterGlass(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	today := time.Now().Format("2006-01-02")
	var water models.WaterIntake

	if err := database.DB.Where("user_id = ? AND date = ?", userID, today).First(&water).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "No water intake record for today"})
		return
	}

	if water.Glasses > 0 {
		water.Glasses--
		water.UpdatedAt = time.Now()
		database.DB.Save(&water)
	}

	c.JSON(http.StatusOK, models.WaterIntakeResponse{
		ID:         water.ID,
		Glasses:    water.Glasses,
		Goal:       water.Goal,
		Date:       water.Date,
		Percentage: water.GetPercentage(),
		Remaining:  water.GetRemaining(),
	})
}

// UpdateWaterGoal updates the daily water goal
func UpdateWaterGoal(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var input struct {
		Goal int `json:"goal" binding:"required,min=1,max=20"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	today := time.Now().Format("2006-01-02")
	var water models.WaterIntake

	if err := database.DB.Where("user_id = ? AND date = ?", userID, today).First(&water).Error; err != nil {
		water = models.WaterIntake{
			UserID:    userID.(uint),
			Glasses:   0,
			Goal:      input.Goal,
			Date:      today,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		database.DB.Create(&water)
	} else {
		water.Goal = input.Goal
		water.UpdatedAt = time.Now()
		database.DB.Save(&water)
	}

	c.JSON(http.StatusOK, models.WaterIntakeResponse{
		ID:         water.ID,
		Glasses:    water.Glasses,
		Goal:       water.Goal,
		Date:       water.Date,
		Percentage: water.GetPercentage(),
		Remaining:  water.GetRemaining(),
	})
}

// GetWaterHistory returns water intake history for past days
func GetWaterHistory(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var history []models.WaterIntake
	
	// Get last 7 days
	database.DB.Where("user_id = ?", userID).
		Order("date DESC").
		Limit(7).
		Find(&history)

	var response []models.WaterIntakeResponse
	for _, water := range history {
		response = append(response, models.WaterIntakeResponse{
			ID:         water.ID,
			Glasses:    water.Glasses,
			Goal:       water.Goal,
			Date:       water.Date,
			Percentage: water.GetPercentage(),
			Remaining:  water.GetRemaining(),
		})
	}

	c.JSON(http.StatusOK, response)
}

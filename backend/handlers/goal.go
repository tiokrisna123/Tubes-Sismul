package handlers

import (
	"health-tracker/database"
	"health-tracker/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// calculateDaysLeft calculates days remaining until deadline
func calculateDaysLeft(deadline string) int {
	if deadline == "" {
		return -1
	}
	
	deadlineTime, err := time.Parse("2006-01-02", deadline)
	if err != nil {
		return -1
	}
	
	days := int(deadlineTime.Sub(time.Now()).Hours() / 24)
	if days < 0 {
		return 0
	}
	return days
}

// GetGoals returns all goals for the user
func GetGoals(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var goals []models.Goal
	if err := database.DB.Where("user_id = ?", userID).Order("created_at DESC").Find(&goals).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch goals"})
		return
	}

	var response []models.GoalResponse
	for _, goal := range goals {
		response = append(response, models.GoalResponse{
			ID:          goal.ID,
			Title:       goal.Title,
			Description: goal.Description,
			Type:        goal.Type,
			Target:      goal.Target,
			Current:     goal.Current,
			Unit:        goal.Unit,
			Deadline:    goal.Deadline,
			IsCompleted: goal.IsCompleted,
			Progress:    goal.GetProgress(),
			DaysLeft:    calculateDaysLeft(goal.Deadline),
		})
	}

	c.JSON(http.StatusOK, response)
}

// CreateGoal creates a new goal
func CreateGoal(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var input struct {
		Title       string  `json:"title" binding:"required"`
		Description string  `json:"description"`
		Type        string  `json:"type" binding:"required"`
		Target      float64 `json:"target" binding:"required"`
		Unit        string  `json:"unit"`
		Deadline    string  `json:"deadline"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	goal := models.Goal{
		UserID:      userID.(uint),
		Title:       input.Title,
		Description: input.Description,
		Type:        input.Type,
		Target:      input.Target,
		Current:     0,
		Unit:        input.Unit,
		Deadline:    input.Deadline,
		IsCompleted: false,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := database.DB.Create(&goal).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create goal"})
		return
	}

	c.JSON(http.StatusCreated, models.GoalResponse{
		ID:          goal.ID,
		Title:       goal.Title,
		Description: goal.Description,
		Type:        goal.Type,
		Target:      goal.Target,
		Current:     goal.Current,
		Unit:        goal.Unit,
		Deadline:    goal.Deadline,
		IsCompleted: goal.IsCompleted,
		Progress:    goal.GetProgress(),
		DaysLeft:    calculateDaysLeft(goal.Deadline),
	})
}

// UpdateGoalProgress updates the current progress of a goal
func UpdateGoalProgress(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	goalID := c.Param("id")
	var goal models.Goal
	if err := database.DB.Where("id = ? AND user_id = ?", goalID, userID).First(&goal).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Goal not found"})
		return
	}

	var input struct {
		Current float64 `json:"current" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	goal.Current = input.Current
	goal.UpdatedAt = time.Now()

	// Check if goal is completed
	if goal.Current >= goal.Target {
		goal.IsCompleted = true
	}

	database.DB.Save(&goal)

	c.JSON(http.StatusOK, models.GoalResponse{
		ID:          goal.ID,
		Title:       goal.Title,
		Description: goal.Description,
		Type:        goal.Type,
		Target:      goal.Target,
		Current:     goal.Current,
		Unit:        goal.Unit,
		Deadline:    goal.Deadline,
		IsCompleted: goal.IsCompleted,
		Progress:    goal.GetProgress(),
		DaysLeft:    calculateDaysLeft(goal.Deadline),
	})
}

// DeleteGoal deletes a goal
func DeleteGoal(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	goalID := c.Param("id")
	var goal models.Goal
	if err := database.DB.Where("id = ? AND user_id = ?", goalID, userID).First(&goal).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Goal not found"})
		return
	}

	database.DB.Delete(&goal)

	c.JSON(http.StatusOK, gin.H{"message": "Goal deleted"})
}

// ToggleGoalComplete toggles the completion status of a goal
func ToggleGoalComplete(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	goalID := c.Param("id")
	var goal models.Goal
	if err := database.DB.Where("id = ? AND user_id = ?", goalID, userID).First(&goal).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Goal not found"})
		return
	}

	goal.IsCompleted = !goal.IsCompleted
	goal.UpdatedAt = time.Now()
	database.DB.Save(&goal)

	c.JSON(http.StatusOK, models.GoalResponse{
		ID:          goal.ID,
		Title:       goal.Title,
		Description: goal.Description,
		Type:        goal.Type,
		Target:      goal.Target,
		Current:     goal.Current,
		Unit:        goal.Unit,
		Deadline:    goal.Deadline,
		IsCompleted: goal.IsCompleted,
		Progress:    goal.GetProgress(),
		DaysLeft:    calculateDaysLeft(goal.Deadline),
	})
}

// GetGoalStats returns summary of goals
func GetGoalStats(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var total, completed, inProgress int64
	
	database.DB.Model(&models.Goal{}).Where("user_id = ?", userID).Count(&total)
	database.DB.Model(&models.Goal{}).Where("user_id = ? AND is_completed = ?", userID, true).Count(&completed)
	inProgress = total - completed

	c.JSON(http.StatusOK, gin.H{
		"total":       total,
		"completed":   completed,
		"in_progress": inProgress,
	})
}

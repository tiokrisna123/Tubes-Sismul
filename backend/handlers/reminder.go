package handlers

import (
	"health-tracker/database"
	"health-tracker/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// GetReminders returns all reminders for the authenticated user
func GetReminders(c *gin.Context) {
	userID := c.GetUint("user_id")

	var reminders []models.Reminder
	result := database.DB.Where("user_id = ?", userID).Order("time ASC").Find(&reminders)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch reminders"})
		return
	}

	// If no reminders exist, create default ones
	if len(reminders) == 0 {
		defaults := models.DefaultReminders()
		for i := range defaults {
			defaults[i].UserID = userID
			database.DB.Create(&defaults[i])
		}
		database.DB.Where("user_id = ?", userID).Order("time ASC").Find(&reminders)
	}

	// Convert to response
	response := make([]models.ReminderResponse, len(reminders))
	for i, r := range reminders {
		response[i] = r.ToResponse()
	}

	c.JSON(http.StatusOK, gin.H{
		"data":    response,
		"message": "Reminders retrieved successfully",
	})
}

// CreateReminder creates a new reminder
func CreateReminder(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req models.CreateReminderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	reminder := models.Reminder{
		UserID:   userID,
		Type:     req.Type,
		Label:    req.Label,
		Time:     req.Time,
		IsActive: true,
	}

	if err := database.DB.Create(&reminder).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create reminder"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"data":    reminder.ToResponse(),
		"message": "Reminder created successfully",
	})
}

// UpdateReminder updates an existing reminder
func UpdateReminder(c *gin.Context) {
	userID := c.GetUint("user_id")
	reminderID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid reminder ID"})
		return
	}

	var reminder models.Reminder
	if err := database.DB.Where("id = ? AND user_id = ?", reminderID, userID).First(&reminder).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Reminder not found"})
		return
	}

	var req models.UpdateReminderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	// Update fields if provided
	if req.Type != "" {
		reminder.Type = req.Type
	}
	if req.Label != "" {
		reminder.Label = req.Label
	}
	if req.Time != "" {
		reminder.Time = req.Time
	}
	if req.IsActive != nil {
		reminder.IsActive = *req.IsActive
	}

	if err := database.DB.Save(&reminder).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update reminder"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":    reminder.ToResponse(),
		"message": "Reminder updated successfully",
	})
}

// DeleteReminder deletes a reminder
func DeleteReminder(c *gin.Context) {
	userID := c.GetUint("user_id")
	reminderID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid reminder ID"})
		return
	}

	result := database.DB.Where("id = ? AND user_id = ?", reminderID, userID).Delete(&models.Reminder{})
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete reminder"})
		return
	}

	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Reminder not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Reminder deleted successfully",
	})
}

// ToggleReminder toggles the active status of a reminder
func ToggleReminder(c *gin.Context) {
	userID := c.GetUint("user_id")
	reminderID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid reminder ID"})
		return
	}

	var reminder models.Reminder
	if err := database.DB.Where("id = ? AND user_id = ?", reminderID, userID).First(&reminder).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Reminder not found"})
		return
	}

	reminder.IsActive = !reminder.IsActive

	if err := database.DB.Save(&reminder).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to toggle reminder"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":    reminder.ToResponse(),
		"message": "Reminder toggled successfully",
	})
}

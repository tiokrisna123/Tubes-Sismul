package handlers

import (
	"net/http"

	"health-tracker/database"
	"health-tracker/models"
	"health-tracker/utils"

	"github.com/gin-gonic/gin"
)

// Register creates a new user account
func Register(c *gin.Context) {
	var req models.RegisterRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request: "+err.Error())
		return
	}

	// Check if email already exists
	var existingUser models.User
	if result := database.DB.Where("email = ?", req.Email).First(&existingUser); result.RowsAffected > 0 {
		utils.ErrorResponse(c, http.StatusConflict, "Email already registered")
		return
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to process password")
		return
	}

	// Create user
	user := models.User{
		Email:    req.Email,
		Password: hashedPassword,
		Name:     req.Name,
	}

	if result := database.DB.Create(&user); result.Error != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to create user")
		return
	}

	// Generate token
	token, err := utils.GenerateToken(user.ID, user.Email)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to generate token")
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "Registration successful", models.LoginResponse{
		Token: token,
		User:  user,
	})
}

// Login authenticates user and returns JWT
func Login(c *gin.Context) {
	var req models.LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request: "+err.Error())
		return
	}

	// Find user by email
	var user models.User
	if result := database.DB.Where("email = ?", req.Email).First(&user); result.Error != nil {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Invalid email or password")
		return
	}

	// Check password
	if !utils.CheckPassword(req.Password, user.Password) {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Invalid email or password")
		return
	}

	// Generate token
	token, err := utils.GenerateToken(user.ID, user.Email)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to generate token")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Login successful", models.LoginResponse{
		Token: token,
		User:  user,
	})
}

// GetCurrentUser returns the authenticated user's profile
func GetCurrentUser(c *gin.Context) {
	userID := c.GetUint("userID")

	var user models.User
	if result := database.DB.First(&user, userID); result.Error != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "User not found")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "User profile retrieved", user)
}

// UpdateProfile updates user profile information
func UpdateProfile(c *gin.Context) {
	userID := c.GetUint("userID")

	var req models.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request: "+err.Error())
		return
	}

	var user models.User
	if result := database.DB.First(&user, userID); result.Error != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "User not found")
		return
	}

	// Update fields
	if req.Name != "" {
		user.Name = req.Name
	}
	if !req.BirthDate.IsZero() {
		user.BirthDate = req.BirthDate
	}
	if req.HeightCm > 0 {
		user.HeightCm = req.HeightCm
	}
	if req.WeightKg > 0 {
		user.WeightKg = req.WeightKg
	}
	if req.ActivityLevel != "" {
		user.ActivityLevel = req.ActivityLevel
	}

	database.DB.Save(&user)

	utils.SuccessResponse(c, http.StatusOK, "Profile updated", user)
}

// ResetPasswordRequest represents the request body for password reset
type ResetPasswordRequest struct {
	Email       string `json:"email" binding:"required,email"`
	NewPassword string `json:"new_password" binding:"required,min=6"`
}

// ResetPassword resets user password by email verification
func ResetPassword(c *gin.Context) {
	var req ResetPasswordRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request: "+err.Error())
		return
	}

	// Find user by email
	var user models.User
	if result := database.DB.Where("email = ?", req.Email).First(&user); result.Error != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Email tidak ditemukan dalam sistem")
		return
	}

	// Hash new password
	hashedPassword, err := utils.HashPassword(req.NewPassword)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to process password")
		return
	}

	// Update password
	user.Password = hashedPassword
	if result := database.DB.Save(&user); result.Error != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to update password")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Password berhasil direset", nil)
}

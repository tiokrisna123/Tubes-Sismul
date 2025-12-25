package handlers

import (
	"net/http"
	"strconv"

	"health-tracker/database"
	"health-tracker/models"
	"health-tracker/utils"

	"github.com/gin-gonic/gin"
)

// InviteFamilyMember sends an invitation to a family member
func InviteFamilyMember(c *gin.Context) {
	userID := c.GetUint("userID")

	var req models.FamilyInviteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request: "+err.Error())
		return
	}

	// Check if member exists
	var memberUser models.User
	if result := database.DB.Where("email = ?", req.MemberEmail).First(&memberUser); result.Error != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "User with this email not found. They need to register first.")
		return
	}

	// Can't invite yourself
	if memberUser.ID == userID {
		utils.ErrorResponse(c, http.StatusBadRequest, "Cannot invite yourself")
		return
	}

	// Check if invitation already exists
	var existing models.FamilyMember
	if result := database.DB.Where("owner_id = ? AND member_user_id = ?", userID, memberUser.ID).First(&existing); result.RowsAffected > 0 {
		utils.ErrorResponse(c, http.StatusConflict, "Invitation already sent to this user")
		return
	}

	// Create invitation
	invitation := models.FamilyMember{
		OwnerID:       userID,
		MemberUserID:  memberUser.ID,
		MemberEmail:   req.MemberEmail,
		Relationship:  req.Relationship,
		Status:        "pending",
		CanViewHealth: true,
	}

	if result := database.DB.Create(&invitation); result.Error != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to send invitation")
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "Invitation sent successfully", invitation)
}

// GetFamilyMembers returns approved family members
func GetFamilyMembers(c *gin.Context) {
	userID := c.GetUint("userID")

	var members []models.FamilyMember
	database.DB.Where("owner_id = ? AND status = ?", userID, "approved").Find(&members)

	// Get user details for each member
	var response []models.FamilyMemberResponse
	for _, m := range members {
		var user models.User
		database.DB.First(&user, m.MemberUserID)
		response = append(response, models.FamilyMemberResponse{
			ID:            m.ID,
			MemberEmail:   m.MemberEmail,
			MemberName:    user.Name,
			Relationship:  m.Relationship,
			Status:        m.Status,
			CanViewHealth: m.CanViewHealth,
			CreatedAt:     m.CreatedAt,
		})
	}

	utils.SuccessResponse(c, http.StatusOK, "Family members retrieved", response)
}

// GetFamilyRequests returns pending invitations for current user
func GetFamilyRequests(c *gin.Context) {
	userID := c.GetUint("userID")

	// Invitations sent to me (where I am the member)
	var received []models.FamilyMember
	database.DB.Where("member_user_id = ? AND status = ?", userID, "pending").Find(&received)

	// Get owner details
	var receivedResponse []map[string]interface{}
	for _, r := range received {
		var owner models.User
		database.DB.First(&owner, r.OwnerID)
		receivedResponse = append(receivedResponse, map[string]interface{}{
			"id":           r.ID,
			"from_email":   owner.Email,
			"from_name":    owner.Name,
			"relationship": r.Relationship,
			"created_at":   r.CreatedAt,
		})
	}

	// Invitations I sent
	var sent []models.FamilyMember
	database.DB.Where("owner_id = ? AND status = ?", userID, "pending").Find(&sent)

	utils.SuccessResponse(c, http.StatusOK, "Family requests retrieved", gin.H{
		"received": receivedResponse,
		"sent":     sent,
	})
}

// ApproveFamilyRequest approves a family invitation
func ApproveFamilyRequest(c *gin.Context) {
	userID := c.GetUint("userID")
	requestID, _ := strconv.ParseUint(c.Param("id"), 10, 32)

	var invitation models.FamilyMember
	if result := database.DB.Where("id = ? AND member_user_id = ? AND status = ?", requestID, userID, "pending").First(&invitation); result.Error != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Invitation not found")
		return
	}

	invitation.Status = "approved"
	database.DB.Save(&invitation)

	// Create reverse relationship so member can also see owner's health
	reverseRelation := models.FamilyMember{
		OwnerID:       userID,
		MemberUserID:  invitation.OwnerID,
		MemberEmail:   "",
		Relationship:  getReverseRelationship(invitation.Relationship),
		Status:        "approved",
		CanViewHealth: true,
	}
	
	// Check if reverse doesn't exist
	var existing models.FamilyMember
	if database.DB.Where("owner_id = ? AND member_user_id = ?", userID, invitation.OwnerID).First(&existing).RowsAffected == 0 {
		database.DB.Create(&reverseRelation)
	}

	utils.SuccessResponse(c, http.StatusOK, "Invitation approved", invitation)
}

// RejectFamilyRequest rejects a family invitation
func RejectFamilyRequest(c *gin.Context) {
	userID := c.GetUint("userID")
	requestID, _ := strconv.ParseUint(c.Param("id"), 10, 32)

	var invitation models.FamilyMember
	if result := database.DB.Where("id = ? AND member_user_id = ? AND status = ?", requestID, userID, "pending").First(&invitation); result.Error != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Invitation not found")
		return
	}

	invitation.Status = "rejected"
	database.DB.Save(&invitation)

	utils.SuccessResponse(c, http.StatusOK, "Invitation rejected", nil)
}

// GetFamilyMemberHealth returns health data of a family member
func GetFamilyMemberHealth(c *gin.Context) {
	userID := c.GetUint("userID")
	familyMemberID, _ := strconv.ParseUint(c.Param("id"), 10, 32)

	// Find family member record by its ID
	var familyMember models.FamilyMember
	if result := database.DB.Where("id = ? AND owner_id = ? AND status = ? AND can_view_health = ?", 
		familyMemberID, userID, "approved", true).First(&familyMember); result.Error != nil {
		utils.ErrorResponse(c, http.StatusForbidden, "You don't have permission to view this member's health")
		return
	}

	// Get the actual member user ID from the family member record
	memberUserID := familyMember.MemberUserID

	// Get member info
	var memberUser models.User
	database.DB.First(&memberUser, memberUserID)

	// Get latest health data
	var latestHealth models.HealthData
	database.DB.Where("user_id = ?", memberUserID).Order("record_date desc").First(&latestHealth)

	// Get recent symptoms
	var recentSymptoms []models.Symptom
	database.DB.Where("user_id = ?", memberUserID).Order("logged_at desc").Limit(5).Find(&recentSymptoms)

	response := models.FamilyHealthView{
		MemberName:     memberUser.Name,
		Relationship:   familyMember.Relationship,
		LatestHealth:   &latestHealth,
		BMICategory:    models.GetBMICategory(latestHealth.BMI),
		RecentSymptoms: recentSymptoms,
	}

	utils.SuccessResponse(c, http.StatusOK, "Family member health retrieved", response)
}

// RemoveFamilyMember removes a family member connection
func RemoveFamilyMember(c *gin.Context) {
	userID := c.GetUint("userID")
	memberID, _ := strconv.ParseUint(c.Param("id"), 10, 32)

	result := database.DB.Where("owner_id = ? AND member_user_id = ?", userID, memberID).Delete(&models.FamilyMember{})
	if result.RowsAffected == 0 {
		utils.ErrorResponse(c, http.StatusNotFound, "Family member not found")
		return
	}

	// Also remove reverse relationship
	database.DB.Where("owner_id = ? AND member_user_id = ?", memberID, userID).Delete(&models.FamilyMember{})

	utils.SuccessResponse(c, http.StatusOK, "Family member removed", nil)
}

func getReverseRelationship(rel string) string {
	switch rel {
	case "parent":
		return "child"
	case "child":
		return "parent"
	case "spouse":
		return "spouse"
	case "sibling":
		return "sibling"
	default:
		return "family"
	}
}

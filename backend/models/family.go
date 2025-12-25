package models

import (
	"time"
)

type FamilyMember struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	OwnerID       uint      `gorm:"not null" json:"owner_id"`
	MemberUserID  uint      `gorm:"not null" json:"member_user_id"`
	MemberEmail   string    `gorm:"not null" json:"member_email"`
	Relationship  string    `json:"relationship"` // parent, child, spouse, sibling, other
	Status        string    `gorm:"default:'pending'" json:"status"` // pending, approved, rejected
	CanViewHealth bool      `gorm:"default:true" json:"can_view_health"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type FamilyInviteRequest struct {
	MemberEmail  string `json:"member_email" binding:"required,email"`
	Relationship string `json:"relationship" binding:"required"`
}

type FamilyMemberResponse struct {
	ID            uint      `json:"id"`
	MemberEmail   string    `json:"member_email"`
	MemberName    string    `json:"member_name"`
	Relationship  string    `json:"relationship"`
	Status        string    `json:"status"`
	CanViewHealth bool      `json:"can_view_health"`
	CreatedAt     time.Time `json:"created_at"`
}

type FamilyHealthView struct {
	MemberName     string       `json:"member_name"`
	Relationship   string       `json:"relationship"`
	LatestHealth   *HealthData  `json:"latest_health"`
	BMICategory    string       `json:"bmi_category"`
	RecentSymptoms []Symptom    `json:"recent_symptoms"`
}

package models

type RegisterUserRequest struct {
	Name           string `json:"name" binding:"required"`
	Email          string `json:"email" binding:"required,email"`
	IsReferral     bool   `json:"isReferral"`
	ReferralUserID string `json:"referralUserId"`
}

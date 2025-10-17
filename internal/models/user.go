package models

type RegisterUserRequestObject struct {
	Name              string `json:"name"`
	Email             string `json:"email"`
	Password          string `json:"password"`
	IsReferral        bool   `json:"is_referral"`
	ReferralUserEmail string `json:"referral_user_email"`
}

type LoginUserRequestObject struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

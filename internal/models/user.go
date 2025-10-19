package models

type RequestBodyRegisterUser struct {
	Name              string `json:"name"`
	Email             string `json:"email"`
	Password          string `json:"password"`
	IsReferral        bool   `json:"is_referral"`
	ReferralUserEmail string `json:"referral_user_email"`
}

type RequestBodyLoginUser struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

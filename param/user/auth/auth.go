package userauthparam

type UserSendLoginCodeRequest struct {
	PhoneNumber string `json:"phone_number"`
}
type UserSendLoginCodeResponse struct {
	// IsNewUser is true when no account exists for this phone number yet, so the
	// frontend can show the "کد دعوت" (ref) field / onboarding on the next step.
	IsNewUser bool `json:"is_new_user"`
}

type UserVerifyLoginCodeRequest struct {
	Code        string `json:"code"`
	PhoneNumber string `json:"phone_number"`
	// Ref is the optional referral/invite code the user arrived with
	// (?ref=<code> on the invite link). Only used when this verification
	// results in a brand-new account.
	Ref string `json:"ref"`
	// Name and BirthDate (Gregorian, "YYYY-MM-DD") are required only when this
	// verification creates a brand-new account (is_new_user from send-code);
	// ignored for existing users.
	Name      string `json:"name"`
	BirthDate string `json:"birth_date"`
}

type UserVerifyLoginCodeResponse struct {
	Token string `json:"token"`
	// IsNewUser is true when this verification just created the account.
	IsNewUser bool `json:"is_new_user"`
}

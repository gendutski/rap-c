package payloadentity

type AttemptLoginPayload struct {
	Email      string `json:"email" form:"email" validate:"required,email"`
	Password   string `json:"password" form:"password" validate:"required"`
	RememberMe bool   `json:"rememberMe" form:"rememberMe"`
}

type RenewPasswordPayload struct {
	Password        string `json:"password" form:"password" validate:"required"`
	ConfirmPassword string `json:"confirmPassword" form:"confirmPassword" validate:"required,eqfield=Password"`
}

type RequestResetPayload struct {
	Email string `json:"email" form:"email" validate:"required,email"`
}

type ValidateResetTokenPayload struct {
	Email string `json:"email" form:"email" query:"email" validate:"required,email"`
	Token string `json:"token" form:"token" query:"token" validate:"required"`
}

type ResetPasswordPayload struct {
	Email           string `json:"email" form:"email" query:"email" validate:"required,email"`
	Token           string `json:"token" form:"token" query:"token" validate:"required"`
	Password        string `json:"password" form:"password" validate:"required"`
	ConfirmPassword string `json:"confirmPassword" form:"confirmPassword" validate:"required,eqfield=Password"`
}

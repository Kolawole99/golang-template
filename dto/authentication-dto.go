package dto

// LoginDTO is used to handle the payload from clients trying to log an existing user in with a POST request
type LoginDTO struct {
	Email    string `json:"email" form:"email" validate:"required,email"`
	Password string `json:"password" form:"password" validate:"required,min=6,max=16"`
}

// LoginOrRegisterResponseDTO is used to handle the response to clients trying to login with a POST request
type LoginOrRegisterResponseDTO struct {
	ID    uint64 `json:"id"`
	Token string `json:"token"`
}

// RegisterDTO is used to handle the payload from clients trying to self-register with a POST request
type RegisterDTO struct {
	Name     string `json:"name" form:"name" validate:"required,min=2,max=30"`
	Email    string `json:"email" form:"email" validate:"required,email"`
	Password string `json:"password" form:"password" validate:"required,min=6,max=16"`
}

// ForgotPasswordDTO is used to handle the payload from clients trying change their password with a POST request
type ForgotPasswordDTO struct {
	Email string `json:"email" form:"email" validate:"required,email"`
}

// ResetPasswordDTO is used to handle the payload from clients trying update their password with a POST request
type ResetPasswordDTO struct {
	Email      string `json:"email" form:"email" binding:"required,email"`
	Token      string `json:"token" form:"token" binding:"required"`
	Password   string `json:"password" form:"password" binding:"required,min=6,max=16"`
	RePassword string `json:"rePassword" form:"rePassword" binding:"required,min=6,max=16"`
}

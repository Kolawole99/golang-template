package controller

import (
	"golang-api/dto"
	"golang-api/helper"
	"golang-api/service"

	"errors"
	"net/http"
)

// AuthController contracts what the Authentication controller can do
type AuthController interface {
	Login(w http.ResponseWriter, r *http.Request)
	Register(w http.ResponseWriter, r *http.Request)
	ForgotPassword(w http.ResponseWriter, r *http.Request)
	ResetPassword(w http.ResponseWriter, r *http.Request)
}

// authController keeps the communication service here
type authController struct {
	authService service.AuthService
}

// NewAuthController creates a new instance of AuthController
func NewAuthController(authService service.AuthService) AuthController {
	return &authController{
		authService: authService,
	}
}

// Login godoc
// @Summary     This Authenticates a user
// @Description This Authenticates an already registered user. It validates the email, password and generates a JWT token along side
// @Tags        Authentication
// @Accept      json
// @Produce     json
// @Router      /auth/login [post]
func (c *authController) Login(w http.ResponseWriter, r *http.Request) {
	var loginDTO dto.LoginDTO

	DTOInJSON := helper.ReadJSON(w, r, &loginDTO)
	if DTOInJSON != nil {
		response := helper.BuildErrorResponse("Cannot process request", DTOInJSON)

		helper.WriteJSON(w, http.StatusConflict, response)

		return
	}

	errorInDTO := helper.ValidateJSON(loginDTO)
	if errorInDTO != nil {
		response := helper.BuildErrorResponse("Failed Validation while Processing Request", errorInDTO)

		helper.WriteJSON(w, http.StatusPreconditionFailed, response)

		return
	}

	userRes := c.authService.LogUserIn(loginDTO.Email, loginDTO.Password)
	if userRes == false {
		response := helper.BuildErrorResponse("Invalid credentials cannot login", errors.New("invalid Username or Password"))

		helper.WriteJSON(w, http.StatusBadRequest, response)

		return
	}

	response := helper.BuildSuccessResponse("Authenticated", userRes)
	helper.WriteJSON(w, http.StatusOK, response)
}

// Register godoc
// @Summary     This creates a new user
// @Description This method handles a User creation payload and creates a not already registered user. It validates the email, password and generates a JWT token along side
// @Tags        Authentication
// @Accept      json
// @Produce     json
// @Router      /auth/register [post]
func (c *authController) Register(w http.ResponseWriter, r *http.Request) {
	var registerDTO dto.RegisterDTO

	DTOInJSON := helper.ReadJSON(w, r, &registerDTO)
	if DTOInJSON != nil {
		response := helper.BuildErrorResponse("Cannot process request", DTOInJSON)

		helper.WriteJSON(w, http.StatusConflict, response)

		return
	}

	errorInDTO := helper.ValidateJSON(registerDTO)
	if errorInDTO != nil {
		response := helper.BuildErrorResponse("Failed Validation while Processing Request", errorInDTO)

		helper.WriteJSON(w, http.StatusPreconditionFailed, response)

		return
	}

	if c.authService.IsDuplicateEmail(registerDTO.Email) {
		response := helper.BuildErrorResponse("Failed to Process Request", errors.New("duplicate email"))

		helper.WriteJSON(w, http.StatusPreconditionFailed, response)

		return
	}

	createdUser := c.authService.CreateUser(registerDTO)

	response := helper.BuildSuccessResponse("Registration successful", createdUser)
	helper.WriteJSON(w, http.StatusCreated, response)
}

// ForgotPassword godoc
// @Summary     This resets a user password
// @Description This disables an already registered user password. It generates a token and send to the email of the account
// @Tags        Authentication
// @Accept      json
// @Produce     json
// @Router      /auth/forgot-password [post]
func (c *authController) ForgotPassword(w http.ResponseWriter, r *http.Request) {
	var forgotPasswordDTO dto.ForgotPasswordDTO

	DTOInJSON := helper.ReadJSON(w, r, &forgotPasswordDTO)
	if DTOInJSON != nil {
		response := helper.BuildErrorResponse("Cannot process request", DTOInJSON)

		helper.WriteJSON(w, http.StatusConflict, response)

		return
	}

	errorInDTO := helper.ValidateJSON(forgotPasswordDTO)
	if errorInDTO != nil {
		response := helper.BuildErrorResponse("Failed Validation while Processing Request", errorInDTO)

		helper.WriteJSON(w, http.StatusPreconditionFailed, response)

		return
	}

	// Generates temporary token for user password reset
	resetSuccessful := c.authService.ForgotPassword(forgotPasswordDTO)
	if !resetSuccessful {
		response := helper.BuildErrorResponse("Password reset failed", errors.New("password reset failed check your email"))

		helper.WriteJSON(w, http.StatusBadRequest, response)

		return
	}

	response := helper.BuildSuccessResponse("Password reset email sent", map[string]string{"email": forgotPasswordDTO.Email})
	helper.WriteJSON(w, http.StatusOK, response)
}

// ResetPassword godoc
// @Summary     This changes a user password
// @Description This updates an already registered user password. It generates a token and logs them into the account afresh
// @Tags        Authentication
// @Accept      json
// @Produce     json
// @Router      /auth/reset-password [post]
func (c *authController) ResetPassword(w http.ResponseWriter, r *http.Request) {
	var resetPasswordDTO dto.ResetPasswordDTO

	DTOInJSON := helper.ReadJSON(w, r, &resetPasswordDTO)
	if DTOInJSON != nil {
		response := helper.BuildErrorResponse("Cannot process request", DTOInJSON)

		helper.WriteJSON(w, http.StatusConflict, response)

		return
	}

	errorInDTO := helper.ValidateJSON(resetPasswordDTO)
	if errorInDTO != nil {
		response := helper.BuildErrorResponse("Failed Validation while Processing Request", errorInDTO)

		helper.WriteJSON(w, http.StatusPreconditionFailed, response)

		return
	}

	if resetPasswordDTO.Password != resetPasswordDTO.RePassword {
		response := helper.BuildErrorResponse("Failed to Process Request", errors.New("passwords do not match"))

		helper.WriteJSON(w, http.StatusPreconditionFailed, response)

		return
	}

	if !c.authService.IsDuplicateEmail(resetPasswordDTO.Email) {
		response := helper.BuildErrorResponse("Failed to Process Request", errors.New("user does not exist"))

		helper.WriteJSON(w, http.StatusPreconditionFailed, response)

		return
	}

	resetSuccessful, err := c.authService.ResetPassword(resetPasswordDTO)
	if err != nil {
		response := helper.BuildErrorResponse("Password reset failed", err)

		helper.WriteJSON(w, http.StatusPreconditionFailed, response)

		return
	}

	response := helper.BuildSuccessResponse("Password reset successful, proceed to Login", resetSuccessful)
	helper.WriteJSON(w, http.StatusOK, response)
}

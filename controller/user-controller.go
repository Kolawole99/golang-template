package controller

import (
	"golang-api/dto"
	"golang-api/helper"
	"golang-api/service"

	"net/http"
	"strconv"
)

// UserController contracts what the User controller can do
type UserController interface {
	UpdateUser(w http.ResponseWriter, r *http.Request)
	Profile(w http.ResponseWriter, r *http.Request)
}

// userController keeps the communication service here
type userController struct {
	userService service.UserService
}

// NewUserController creates a new instance of UserController
func NewUserController(userService service.UserService) UserController {
	return &userController{
		userService: userService,
	}
}

// UpdateUser godoc
// @Summary     This updates a user
// @Description This updates an already existing user. It validates the details of the user matches the authenticated user and then updates the user
// @Tags        Users
// @Accept      json
// @Produce     json
// @Router      /users/update [put]
func (c *userController) UpdateUser(w http.ResponseWriter, r *http.Request) {
	var userUpdateDTO dto.UserUpdateDTO

	DTOInJSON := helper.ReadJSON(w, r, &userUpdateDTO)
	if DTOInJSON != nil {
		response := helper.BuildErrorResponse("Cannot process request", DTOInJSON)

		helper.WriteJSON(w, http.StatusConflict, response)

		return
	}

	errorInDTO := helper.ValidateJSON(userUpdateDTO)
	if errorInDTO != nil {
		response := helper.BuildErrorResponse("Failed Validation while Processing Request", errorInDTO)

		helper.WriteJSON(w, http.StatusPreconditionFailed, response)

		return
	}

	id, err := strconv.ParseUint(r.Header.Get(helper.CURRENT_USER), 10, 64)
	if err != nil {
		response := helper.BuildErrorResponse("Failed to Process Request", err)

		helper.WriteJSON(w, http.StatusInternalServerError, response)

		return
	}

	userUpdateDTO.ID = id
	updatedUser := c.userService.UpdateUser(userUpdateDTO)
	response := helper.BuildSuccessResponse("Ok", updatedUser)

	helper.WriteJSON(w, http.StatusOK, response)
}

// Profile godoc
// @Summary     This returns a user profile
// @Description This returns an already existing user. It validates the authenticated user matches the token and returns their details
// @Tags        Users
// @Accept      json
// @Produce     json
// @Router      /users/profile [get]
func (c *userController) Profile(w http.ResponseWriter, r *http.Request) {
	authenticatedUserId := r.Header.Get(helper.CURRENT_USER)

	profile := c.userService.Profile(authenticatedUserId)

	response := helper.BuildSuccessResponse("Ok", profile)

	helper.WriteJSON(w, http.StatusOK, response)
}

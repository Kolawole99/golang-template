package service

import (
	"golang-api/dto"
	"golang-api/helper"
	"golang-api/model"

	"errors"
	"log"
	"strconv"

	"github.com/mashingan/smapping"
)

// AuthService is a contract of what the Authentication service can do
type AuthService interface {
	LogUserIn(email string, password string) interface{}
	CreateUser(user dto.RegisterDTO) model.User
	ForgotPassword(user dto.ForgotPasswordDTO) bool
	ResetPassword(user dto.ResetPasswordDTO) (*model.User, error)
	IsDuplicateEmail(email string) bool
}

type authService struct {
	userRepository model.UserRepository
	jwtService     JWTService
}

// NewAuthService creates a new AuthService
func NewAuthService(userRepository model.UserRepository, jwtService JWTService) AuthService {
	return &authService{
		userRepository: userRepository,
		jwtService:     jwtService,
	}
}

func (service *authService) LogUserIn(email string, password string) interface{} {
	user := service.userRepository.FindByEmail(email)
	if user.Email != email {
		return false
	}

	isValid := helper.ComparePassword(user.Password, []byte(password))
	if !isValid {
		return false
	}

	generatedToken := service.jwtService.GenerateToken(strconv.FormatUint(user.ID, 10))
	user.Token = generatedToken

	return user
}

func (service *authService) CreateUser(user dto.RegisterDTO) model.User {
	userToCreate := model.User{}

	err := smapping.FillStruct(&userToCreate, smapping.MapFields(&user))
	if err != nil {
		log.Fatalf("Failed to map User DTO to model with error %v", err)
	}

	userToCreate.Password = helper.HashPassword([]byte(userToCreate.Password))

	createdUser := service.userRepository.InsertUser(userToCreate)

	generatedToken := service.jwtService.GenerateToken(strconv.FormatUint(createdUser.ID, 10))
	createdUser.Token = generatedToken

	return createdUser
}

func (service *authService) ForgotPassword(user dto.ForgotPasswordDTO) bool {
	userToUpdate := model.User{}

	err := smapping.FillStruct(&userToUpdate, smapping.MapFields(&user))
	if err != nil {
		log.Fatalf("Failed to map User Forgot Password DTO to model with error %v", err)

		return false
	}

	res := service.userRepository.IsDuplicateEmail(userToUpdate.Email)
	if res.RowsAffected < 1 {
		return false
	}

	resetToken := helper.GenerateToken()
	userToUpdate.PasswordResetToken = helper.HashPassword([]byte(resetToken))
	userToUpdate.IsResettingPassword = true

	createdUser := service.userRepository.UpdateUser(userToUpdate)
	if createdUser.Error != nil {
		return false
	}

	// TASK
	// Sends email of password reset to user
	// resetDetails := map[string]string{
	// 	"email":         userToUpdate.Email,
	// "name": userToUpdate.Name
	// 	"unHashedToken": resetToken,
	// }

	// fmt.Println(resetDetails, "Pwd forgot-password email here")

	return true
}

func (service *authService) ResetPassword(resetData dto.ResetPasswordDTO) (*model.User, error) {
	user := service.userRepository.FindByEmail(resetData.Email)

	if !user.IsResettingPassword {
		return nil, errors.New("no password reset started")
	}

	isValid := helper.ComparePassword(user.PasswordResetToken, []byte(resetData.Token))
	if !isValid {
		return nil, errors.New("invalid Token")
	}

	user.Password = helper.HashPassword([]byte(resetData.Password))
	user.IsResettingPassword = false
	user.PasswordResetToken = ""

	updatedUser := service.userRepository.UpdateUser(user)
	if updatedUser.Error != nil {
		log.Fatalf("Error updating user %v", updatedUser.Error)

		return nil, updatedUser.Error
	}

	return &model.User{}, nil
}

func (service *authService) IsDuplicateEmail(email string) bool {
	res := service.userRepository.IsDuplicateEmail(email)

	return res.Error == nil
}

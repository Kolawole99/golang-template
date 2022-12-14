package service

import (
	"golang-api/dto"
	"golang-api/helper"
	"golang-api/model"

	"fmt"
	"log"

	"github.com/mashingan/smapping"
)

// UserService is a contract of what this service can do
type UserService interface {
	UpdateUser(user dto.UserUpdateDTO) model.User
	Profile(userId string) model.User
}

type userService struct {
	userRepository model.UserRepository
}

// NewUserService creates a new UserService
func NewUserService(userRepo model.UserRepository) UserService {
	return &userService{
		userRepository: userRepo,
	}
}

func (u *userService) UpdateUser(user dto.UserUpdateDTO) model.User {
	userToUpdate := model.User{}

	err := smapping.FillStruct(&userToUpdate, smapping.MapFields(&user))
	if err != nil {
		log.Fatalf("Error in updating User")
	}

	if userToUpdate.Password != "" {
		userToUpdate.Password = helper.HashPassword([]byte(userToUpdate.Password))
	} else {
		tempUser := u.userRepository.DetailsUser(fmt.Sprintf("%v", userToUpdate.ID))

		userToUpdate.Password = tempUser.Password
	}

	updatedUser := u.userRepository.InsertUser(userToUpdate)

	return updatedUser
}

func (u *userService) Profile(userId string) model.User {
	return u.userRepository.ProfileUser(userId)
}

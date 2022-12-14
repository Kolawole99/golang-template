package model

import (
	"gorm.io/gorm"
)

// User represents a user and the users table in the database
type User struct {
	ID                  uint64 `gorm:"primary_key:auto_increment" json:"id"`
	Name                string `gorm:"type:varchar(255)" json:"name"`
	Email               string `gorm:"uniqueIndex;type:varchar(255)" json:"email"`
	Password            string `gorm:"->;<-;not null" json:"-"`
	Token               string `gorm:"" json:"token,omitempty"`
	IsResettingPassword bool   `gorm:"default:false" json:"-"`
	PasswordResetToken  string `gorm:"type:varchar(255)" json:"-"`
}

// UserRepository is a contract of what userRepository can do to the DB
type UserRepository interface {
	InsertUser(user User) User
	UpdateUser(user User) (tx *gorm.DB)
	IsDuplicateEmail(email string) (tx *gorm.DB)
	FindByEmail(email string) User
	DetailsUser(userId string) User
	ProfileUser(userId string) User
}

type userConnection struct {
	connection *gorm.DB
}

// NewUserRepository creates a new instance of UserRepository
func NewUserRepository(db *gorm.DB) UserRepository {
	return &userConnection{
		connection: db,
	}
}

// Handles Creating Operations. Would create if row does not already exist
func (c *userConnection) InsertUser(user User) User {
	c.connection.Save(&user)

	return user
}

// Handle update column, would update by email index
func (c *userConnection) UpdateUser(user User) (tx *gorm.DB) {
	return c.connection.Where("email = ?", user.Email).Updates(&user)
}

func (c *userConnection) IsDuplicateEmail(email string) (tx *gorm.DB) {
	var user User

	return c.connection.Where("email = ?", email).Take(&user)
}

// Handle lookup operation. Searches by email indexed field
func (c *userConnection) FindByEmail(email string) User {
	var user User

	c.connection.Where("email = ?", email).Take(&user)

	return user
}

// Handle lookup operation. Searches by userId primary key field
func (c *userConnection) DetailsUser(userId string) User {
	var tempUser User

	c.connection.Find(&tempUser, userId)

	return tempUser
}

func (c *userConnection) ProfileUser(userId string) User {
	var user User

	c.connection.Find(&user, userId).Take(&user)

	return user
}

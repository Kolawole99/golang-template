package dto

// UserUpdateDTO is used to handle the payload from clients trying to update an existing user with a PUT request
type UserUpdateDTO struct {
	ID       uint64 `json:"id" form:"id" binding:"required"`
	Name     string `json:"name" form:"name" binding:"required"`
	Email    string `json:"email" form:"email" binding:"required"`
	Password string `json:"password,omitempty" form:"password,omitempty"`
}

// UserCreateDTO is used to handle the payload from clients trying to create a new user with a POST request
type UserCreateDTO struct {
	Name     string `json:"name" form:"name" binding:"required"`
	Email    string `json:"email" form:"email" binding:"required,email"`
	Password string `json:"password,omitempty" form:"password,omitempty" binding:"required"`
}

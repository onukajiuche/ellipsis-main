package model

import "time"

type User struct {
	ID        string    `json:"id,omitempty" gorm:"column:id;index;unique;not null;type:varchar(50)"`
	Firstname string    `json:"firstname,omitempty" gorm:"column:firstname;not null;type:varchar(100)" validate:"required"`
	Lastname  string    `json:"lastname,omitempty" gorm:"column:lastname;type:varchar(100)"`
	Email     string    `json:"email,omitempty" gorm:"column:email;index;unique;not null;type:varchar(100)" validate:"email,required"`
	Password  string    `json:"password,omitempty" gorm:"column:password;type:varchar(100);not null" validate:"required,min=8"`
	Role      int       `json:"role,omitempty" gorm:"column:role;not null;type:smallint"`
	IsLocked  bool      `json:"is_locked,omitempty" gorm:"column:is_locked"`
	Salt      string    `json:"salt,omitempty" gorm:"column:salt;not null;type:varchar(50)"`
	CreatedAt time.Time `json:"created_at" gorm:"column:created_at;index"`
	UpdatedAt time.Time `json:"updated_at" gorm:"column:updated_at"`
	Urls      []URL     `json:"-" gorm:"foreignKey:user_id" swaggerignore:"true"`
}

type UserLogin struct {
	Email    string `json:"email,omitempty"`
	Password string `json:"password,omitempty"`
}

type LoginResponse struct {
	Token string `json:"token,omitempty"`
	User  *User  `json:"user"`
}

type ResetPassword struct {
	Password string `json:"password,omitempty"`
	Salt     string `json:"-" swaggerignore:"true"`
}

type ForgotPassword struct {
	Email string `json:"email,omitempty"`
}

type ContextInfo struct {
	ID    string
	Role  int
	Email string
}

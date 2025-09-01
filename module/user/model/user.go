package userModel

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ID uuid.UUID `gorm:"type:uuid"`
	Name string `json:"name"`
	Email string `json:"description"`
	PasswordHash string `json:"password"`
}

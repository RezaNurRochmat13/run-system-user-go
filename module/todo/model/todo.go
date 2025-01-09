package todoModel

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Todo struct {
	gorm.Model
	ID uuid.UUID `gorm:"type:uuid"`
	Title string `json:"title"`
	Description string `json:"description"`
	Status string `json:"status"`
	DueDate string `json:"due_date"`
}
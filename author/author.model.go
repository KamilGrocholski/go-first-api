package author

import (
	"time"

	"gorm.io/gorm"
)

type Author struct {
	gorm.Model
	FullName    string     `json:"fullName" binding:"required"`
	BirthDate   *time.Time `json:"birthDay"`
	ImageUrl    string     `json:"imageUrl"`
	Description string     `json:"description"`
}

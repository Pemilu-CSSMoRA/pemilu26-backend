package entity

import (
	"github.com/Pemilu-CSSMoRA/pemilu26-backend/internal/utils/password"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRole string

const (
	RoleAdmin   UserRole = "admin"
	RoleBawaslu UserRole = "bawaslu"
	RoleUser    UserRole = "user"
)

type User struct {
	ID       uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	Name     string    `json:"name"`
	NIA      string    `json:"nia" gorm:"unique"`
	Angkatan string    `json:"angkatan"`
	Password string    `json:"password"`
	Role     UserRole  `json:"role" gorm:"default:user"`
	Timestamp
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var err error
	// u.ID = uuid.New()
	u.Password, err = password.HashPassword(u.Password)
	if err != nil {
		return err
	}
	return nil
}

func (User) TableName() string {
	return "users"
}

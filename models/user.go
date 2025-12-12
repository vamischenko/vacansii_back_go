package models

import (
	"time"

	"golang.org/x/crypto/bcrypt"
)

// User модель пользователя
// Представляет пользователя системы с безопасной аутентификацией
type User struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	Username     string    `gorm:"type:varchar(255);uniqueIndex;not null" json:"username" binding:"required"`
	Email        string    `gorm:"type:varchar(255);uniqueIndex;not null" json:"email" binding:"required,email"`
	PasswordHash string    `gorm:"type:varchar(255);not null" json:"-"`
	AuthKey      string    `gorm:"type:varchar(32)" json:"-"`
	AccessToken  string    `gorm:"type:varchar(64);uniqueIndex" json:"access_token,omitempty"`
	Status       int       `gorm:"default:10;not null" json:"status"`
	CreatedAt    time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt    time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

const (
	// StatusDeleted статус удаленного пользователя
	StatusDeleted = 0
	// StatusActive статус активного пользователя
	StatusActive = 10
)

// TableName указывает имя таблицы для модели User
func (User) TableName() string {
	return "user"
}

// SetPassword устанавливает хешированный пароль
func (u *User) SetPassword(password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.PasswordHash = string(hashedPassword)
	return nil
}

// ValidatePassword проверяет пароль пользователя
func (u *User) ValidatePassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password))
	return err == nil
}

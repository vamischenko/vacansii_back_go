package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"
)

// JSONB тип для хранения JSON данных в базе
type JSONB map[string]interface{}

// Value преобразует JSONB в значение для базы данных
func (j JSONB) Value() (driver.Value, error) {
	if j == nil {
		return nil, nil
	}
	return json.Marshal(j)
}

// Scan преобразует значение из базы данных в JSONB
func (j *JSONB) Scan(value interface{}) error {
	if value == nil {
		*j = nil
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("failed to unmarshal JSONB value")
	}

	result := make(map[string]interface{})
	err := json.Unmarshal(bytes, &result)
	*j = result
	return err
}

// Vacancy модель вакансии
// Представляет сущность вакансии с основными полями и дополнительными данными в формате JSON
type Vacancy struct {
	ID               uint      `gorm:"primaryKey" json:"id"`
	Title            string    `gorm:"type:varchar(255);not null" json:"title" binding:"required,max=255"`
	Description      string    `gorm:"type:text;not null" json:"description" binding:"required"`
	Salary           int       `gorm:"not null" json:"salary" binding:"required,min=0"`
	AdditionalFields JSONB     `gorm:"type:json" json:"additional_fields,omitempty"`
	CreatedAt        time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt        time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

// TableName указывает имя таблицы для модели Vacancy
func (Vacancy) TableName() string {
	return "vacancy"
}

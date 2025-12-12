package repositories

import (
	"fmt"
	"strings"
	"vakansii-back-go/models"

	"gorm.io/gorm"
)

// VacancyRepository интерфейс для работы с вакансиями
type VacancyRepository interface {
	FindByID(id uint) (*models.Vacancy, error)
	FindAll(page int, sortBy string, sortOrder string) ([]models.Vacancy, int64, error)
	Save(vacancy *models.Vacancy) error
	Update(vacancy *models.Vacancy) error
	Delete(id uint) error
	GetTotalCount() (int64, error)
	Search(query string, page int, sortOrder string) ([]models.Vacancy, int64, error)
}

// vacancyRepository реализация VacancyRepository
type vacancyRepository struct {
	db *gorm.DB
}

const PageSize = 10

// NewVacancyRepository создает новый экземпляр репозитория вакансий
func NewVacancyRepository(db *gorm.DB) VacancyRepository {
	return &vacancyRepository{db: db}
}

// FindByID находит вакансию по ID
func (r *vacancyRepository) FindByID(id uint) (*models.Vacancy, error) {
	var vacancy models.Vacancy
	if err := r.db.First(&vacancy, id).Error; err != nil {
		return nil, err
	}
	return &vacancy, nil
}

// FindAll получает все вакансии с пагинацией и сортировкой
func (r *vacancyRepository) FindAll(page int, sortBy string, sortOrder string) ([]models.Vacancy, int64, error) {
	var vacancies []models.Vacancy
	var total int64

	// Считаем общее количество
	if err := r.db.Model(&models.Vacancy{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Валидация полей сортировки
	allowedSortFields := map[string]bool{
		"salary":     true,
		"created_at": true,
	}
	if !allowedSortFields[sortBy] {
		sortBy = "created_at"
	}

	// Валидация порядка сортировки
	if sortOrder != "asc" && sortOrder != "desc" {
		sortOrder = "desc"
	}

	offset := (page - 1) * PageSize
	orderClause := fmt.Sprintf("%s %s", sortBy, sortOrder)

	if err := r.db.Order(orderClause).Limit(PageSize).Offset(offset).Find(&vacancies).Error; err != nil {
		return nil, 0, err
	}

	return vacancies, total, nil
}

// Save сохраняет новую вакансию
func (r *vacancyRepository) Save(vacancy *models.Vacancy) error {
	return r.db.Create(vacancy).Error
}

// Update обновляет существующую вакансию
func (r *vacancyRepository) Update(vacancy *models.Vacancy) error {
	return r.db.Save(vacancy).Error
}

// Delete удаляет вакансию по ID
func (r *vacancyRepository) Delete(id uint) error {
	result := r.db.Delete(&models.Vacancy{}, id)
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return result.Error
}

// GetTotalCount возвращает общее количество вакансий
func (r *vacancyRepository) GetTotalCount() (int64, error) {
	var count int64
	if err := r.db.Model(&models.Vacancy{}).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// Search выполняет полнотекстовый поиск по вакансиям
func (r *vacancyRepository) Search(query string, page int, sortOrder string) ([]models.Vacancy, int64, error) {
	var vacancies []models.Vacancy
	var total int64

	query = strings.TrimSpace(query)

	// Базовый запрос с FULLTEXT поиском для MySQL
	db := r.db.Model(&models.Vacancy{}).
		Where("MATCH(title, description) AGAINST (? IN NATURAL LANGUAGE MODE)", query)

	// Считаем общее количество результатов
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Применяем сортировку
	offset := (page - 1) * PageSize

	if sortOrder == "relevance" {
		// Сортировка по релевантности (для MySQL)
		db = db.Select("*, MATCH(title, description) AGAINST (? IN NATURAL LANGUAGE MODE) as relevance_score", query).
			Order("relevance_score DESC")
	} else if sortOrder == "asc" {
		db = db.Order("created_at ASC")
	} else {
		db = db.Order("created_at DESC")
	}

	// Применяем пагинацию
	if err := db.Limit(PageSize).Offset(offset).Find(&vacancies).Error; err != nil {
		return nil, 0, err
	}

	return vacancies, total, nil
}

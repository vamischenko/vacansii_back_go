package services

import (
	"fmt"
	"math"
	"strings"
	"vakansii-back-go/models"
	"vakansii-back-go/repositories"

	"gorm.io/gorm"
)

// VacancyService интерфейс сервиса вакансий
type VacancyService interface {
	GetVacancyList(page int, sortBy string, sortOrder string) (map[string]interface{}, error)
	GetVacancyByID(id uint, fields []string) (interface{}, error)
	CreateVacancy(data map[string]interface{}) (map[string]interface{}, error)
	UpdateVacancy(id uint, data map[string]interface{}) (map[string]interface{}, error)
	DeleteVacancy(id uint) (map[string]interface{}, error)
	SearchVacancies(query string, page int, sortOrder string) (map[string]interface{}, error)
}

// vacancyService реализация VacancyService
type vacancyService struct {
	repo repositories.VacancyRepository
}

// NewVacancyService создает новый экземпляр сервиса вакансий
func NewVacancyService(repo repositories.VacancyRepository) VacancyService {
	return &vacancyService{repo: repo}
}

// GetVacancyList получает список вакансий с пагинацией
func (s *vacancyService) GetVacancyList(page int, sortBy string, sortOrder string) (map[string]interface{}, error) {
	vacancies, total, err := s.repo.FindAll(page, sortBy, sortOrder)
	if err != nil {
		return nil, err
	}

	pageCount := int(math.Ceil(float64(total) / float64(repositories.PageSize)))

	return map[string]interface{}{
		"data": vacancies,
		"pagination": map[string]interface{}{
			"total":     total,
			"page":      page,
			"pageSize":  repositories.PageSize,
			"pageCount": pageCount,
		},
	}, nil
}

// GetVacancyByID получает вакансию по ID с возможностью выбора полей
func (s *vacancyService) GetVacancyByID(id uint, fields []string) (interface{}, error) {
	vacancy, err := s.repo.FindByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}

	// Если поля не указаны, возвращаем всю вакансию
	if len(fields) == 0 {
		return vacancy, nil
	}

	// Фильтруем поля
	result := make(map[string]interface{})
	result["id"] = vacancy.ID

	for _, field := range fields {
		field = strings.TrimSpace(field)
		switch field {
		case "title":
			result["title"] = vacancy.Title
		case "description":
			result["description"] = vacancy.Description
		case "salary":
			result["salary"] = vacancy.Salary
		case "additional_fields":
			result["additional_fields"] = vacancy.AdditionalFields
		case "created_at":
			result["created_at"] = vacancy.CreatedAt
		case "updated_at":
			result["updated_at"] = vacancy.UpdatedAt
		}
	}

	return result, nil
}

// CreateVacancy создает новую вакансию
func (s *vacancyService) CreateVacancy(data map[string]interface{}) (map[string]interface{}, error) {
	vacancy := &models.Vacancy{}

	// Заполняем основные поля
	if title, ok := data["title"].(string); ok {
		vacancy.Title = title
	} else {
		return map[string]interface{}{
			"success": false,
			"message": "Название вакансии обязательно",
		}, nil
	}

	if description, ok := data["description"].(string); ok {
		vacancy.Description = description
	} else {
		return map[string]interface{}{
			"success": false,
			"message": "Описание вакансии обязательно",
		}, nil
	}

	if salary, ok := data["salary"].(float64); ok {
		vacancy.Salary = int(salary)
	} else if salaryInt, ok := data["salary"].(int); ok {
		vacancy.Salary = salaryInt
	} else {
		return map[string]interface{}{
			"success": false,
			"message": "Зарплата обязательна",
		}, nil
	}

	// Дополнительные поля (опционально)
	if additionalFields, ok := data["additional_fields"].(map[string]interface{}); ok {
		vacancy.AdditionalFields = additionalFields
	}

	// Сохраняем в базу
	if err := s.repo.Save(vacancy); err != nil {
		return map[string]interface{}{
			"success": false,
			"message": "Ошибка при создании вакансии",
			"error":   err.Error(),
		}, nil
	}

	return map[string]interface{}{
		"success": true,
		"id":      vacancy.ID,
		"message": "Вакансия успешно создана",
	}, nil
}

// UpdateVacancy обновляет существующую вакансию
func (s *vacancyService) UpdateVacancy(id uint, data map[string]interface{}) (map[string]interface{}, error) {
	vacancy, err := s.repo.FindByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return map[string]interface{}{
				"success": false,
				"message": "Вакансия не найдена",
			}, nil
		}
		return nil, err
	}

	// Обновляем поля если они присутствуют
	if title, ok := data["title"].(string); ok {
		vacancy.Title = title
	}

	if description, ok := data["description"].(string); ok {
		vacancy.Description = description
	}

	if salary, ok := data["salary"].(float64); ok {
		vacancy.Salary = int(salary)
	} else if salaryInt, ok := data["salary"].(int); ok {
		vacancy.Salary = salaryInt
	}

	if additionalFields, ok := data["additional_fields"].(map[string]interface{}); ok {
		vacancy.AdditionalFields = additionalFields
	}

	// Сохраняем изменения
	if err := s.repo.Update(vacancy); err != nil {
		return map[string]interface{}{
			"success": false,
			"message": "Ошибка при обновлении вакансии",
			"error":   err.Error(),
		}, nil
	}

	return map[string]interface{}{
		"success": true,
		"message": "Вакансия успешно обновлена",
	}, nil
}

// DeleteVacancy удаляет вакансию
func (s *vacancyService) DeleteVacancy(id uint) (map[string]interface{}, error) {
	err := s.repo.Delete(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return map[string]interface{}{
				"success": false,
				"message": "Вакансия не найдена",
			}, nil
		}
		return nil, err
	}

	return map[string]interface{}{
		"success": true,
		"message": "Вакансия успешно удалена",
	}, nil
}

// SearchVacancies выполняет полнотекстовый поиск вакансий
func (s *vacancyService) SearchVacancies(query string, page int, sortOrder string) (map[string]interface{}, error) {
	query = strings.TrimSpace(query)
	if query == "" {
		return map[string]interface{}{
			"success": false,
			"message": "Поисковый запрос не может быть пустым",
		}, nil
	}

	vacancies, total, err := s.repo.Search(query, page, sortOrder)
	if err != nil {
		return nil, fmt.Errorf("ошибка поиска: %w", err)
	}

	pageCount := int(math.Ceil(float64(total) / float64(repositories.PageSize)))

	return map[string]interface{}{
		"data": vacancies,
		"pagination": map[string]interface{}{
			"total":     total,
			"page":      page,
			"pageSize":  repositories.PageSize,
			"pageCount": pageCount,
		},
		"query": query,
	}, nil
}

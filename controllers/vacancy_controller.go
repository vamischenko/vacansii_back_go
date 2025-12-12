package controllers

import (
	"net/http"
	"strconv"
	"strings"
	"vakansii-back-go/services"

	"github.com/gin-gonic/gin"
)

// VacancyController контроллер для работы с вакансиями
type VacancyController struct {
	service services.VacancyService
}

// NewVacancyController создает новый экземпляр контроллера вакансий
func NewVacancyController(service services.VacancyService) *VacancyController {
	return &VacancyController{service: service}
}

// Index получает список вакансий с пагинацией
// GET /vacancy
func (vc *VacancyController) Index(c *gin.Context) {
	// Получаем параметры запроса
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	sortBy := c.DefaultQuery("sort", "created_at")
	sortOrder := c.DefaultQuery("order", "desc")

	// Валидация page
	if page < 1 {
		page = 1
	}
	if page > 10000 {
		page = 10000
	}

	result, err := vc.service.GetVacancyList(page, sortBy, sortOrder)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Ошибка при получении списка вакансий",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, result)
}

// View получает конкретную вакансию по ID
// GET /vacancy/:id
func (vc *VacancyController) View(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Неверный ID вакансии",
		})
		return
	}

	// Получаем параметр fields
	fieldsParam := c.Query("fields")
	var fields []string
	if fieldsParam != "" {
		fields = strings.Split(fieldsParam, ",")
		// Ограничиваем количество полей до 10
		if len(fields) > 10 {
			fields = fields[:10]
		}
	}

	result, err := vc.service.GetVacancyByID(uint(id), fields)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Ошибка при получении вакансии",
			"error":   err.Error(),
		})
		return
	}

	if result == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "Вакансия не найдена",
		})
		return
	}

	c.JSON(http.StatusOK, result)
}

// Create создает новую вакансию
// POST /vacancy
func (vc *VacancyController) Create(c *gin.Context) {
	var data map[string]interface{}
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Неверный формат данных",
			"error":   err.Error(),
		})
		return
	}

	result, err := vc.service.CreateVacancy(data)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Ошибка при создании вакансии",
			"error":   err.Error(),
		})
		return
	}

	// Если в результате есть success: false, возвращаем 400
	if success, ok := result["success"].(bool); ok && !success {
		c.JSON(http.StatusBadRequest, result)
		return
	}

	c.JSON(http.StatusCreated, result)
}

// Update обновляет существующую вакансию
// PUT /vacancy/:id
func (vc *VacancyController) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Неверный ID вакансии",
		})
		return
	}

	var data map[string]interface{}
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Неверный формат данных",
			"error":   err.Error(),
		})
		return
	}

	result, err := vc.service.UpdateVacancy(uint(id), data)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Ошибка при обновлении вакансии",
			"error":   err.Error(),
		})
		return
	}

	// Если в результате есть success: false, возвращаем 404
	if success, ok := result["success"].(bool); ok && !success {
		c.JSON(http.StatusNotFound, result)
		return
	}

	c.JSON(http.StatusOK, result)
}

// Delete удаляет вакансию
// DELETE /vacancy/:id
func (vc *VacancyController) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Неверный ID вакансии",
		})
		return
	}

	result, err := vc.service.DeleteVacancy(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Ошибка при удалении вакансии",
			"error":   err.Error(),
		})
		return
	}

	// Если в результате есть success: false, возвращаем 404
	if success, ok := result["success"].(bool); ok && !success {
		c.JSON(http.StatusNotFound, result)
		return
	}

	c.JSON(http.StatusOK, result)
}

// Search выполняет полнотекстовый поиск вакансий
// GET /vacancy/search
func (vc *VacancyController) Search(c *gin.Context) {
	query := strings.TrimSpace(c.Query("q"))
	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Поисковый запрос не может быть пустым",
		})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	if page < 1 {
		page = 1
	}
	if page > 10000 {
		page = 10000
	}

	sortOrder := c.DefaultQuery("sort", "relevance")

	result, err := vc.service.SearchVacancies(query, page, sortOrder)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Ошибка при поиске вакансий",
			"error":   err.Error(),
		})
		return
	}

	// Если в результате есть success: false, возвращаем 400
	if success, ok := result["success"].(bool); ok && !success {
		c.JSON(http.StatusBadRequest, result)
		return
	}

	c.JSON(http.StatusOK, result)
}

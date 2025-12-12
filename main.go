package main

import (
	"fmt"
	"log"
	"vakansii-back-go/config"
	"vakansii-back-go/controllers"
	"vakansii-back-go/middleware"
	"vakansii-back-go/migrations"
	"vakansii-back-go/repositories"
	"vakansii-back-go/services"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	// Загружаем конфигурацию
	cfg := config.Load()

	// Подключаемся к базе данных
	db, err := gorm.Open(mysql.Open(cfg.GetDSN()), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Выполняем миграции
	if err := migrations.Migrate(db); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	// Устанавливаем режим Gin
	gin.SetMode(cfg.Server.Mode)

	// Создаем роутер
	r := gin.Default()

	// Настраиваем CORS
	corsConfig := cors.Config{
		AllowOrigins:     cfg.CORS.AllowOrigins,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
		AllowHeaders:     []string{"*"},
		AllowCredentials: false,
		MaxAge:           86400,
	}
	r.Use(cors.New(corsConfig))

	// Добавляем rate limiter middleware
	r.Use(middleware.RateLimiter(cfg.RateLimit.Requests, cfg.RateLimit.Window))

	// Инициализируем слои
	vacancyRepo := repositories.NewVacancyRepository(db)
	vacancyService := services.NewVacancyService(vacancyRepo)
	vacancyController := controllers.NewVacancyController(vacancyService)

	// Настраиваем роуты
	vacancyGroup := r.Group("/vacancy")
	{
		vacancyGroup.GET("", vacancyController.Index)
		vacancyGroup.GET("/search", vacancyController.Search)
		vacancyGroup.GET("/:id", vacancyController.View)
		vacancyGroup.POST("", vacancyController.Create)
		vacancyGroup.PUT("/:id", vacancyController.Update)
		vacancyGroup.DELETE("/:id", vacancyController.Delete)
	}

	// Запускаем сервер
	addr := fmt.Sprintf(":%s", cfg.Server.Port)
	fmt.Printf("Server starting on %s\n", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

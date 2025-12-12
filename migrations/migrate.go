package migrations

import (
	"fmt"
	"vakansii-back-go/models"

	"gorm.io/gorm"
)

// Migrate выполняет автоматические миграции для всех моделей
func Migrate(db *gorm.DB) error {
	fmt.Println("Running database migrations...")

	// Автоматическая миграция для таблиц
	err := db.AutoMigrate(
		&models.Vacancy{},
		&models.User{},
	)

	if err != nil {
		return fmt.Errorf("failed to migrate database: %w", err)
	}

	// Создание FULLTEXT индексов для поиска (для MySQL)
	if err := createFullTextIndexes(db); err != nil {
		fmt.Printf("Warning: failed to create fulltext indexes: %v\n", err)
	}

	fmt.Println("Database migrations completed successfully")
	return nil
}

// createFullTextIndexes создает FULLTEXT индексы для таблицы vacancy
func createFullTextIndexes(db *gorm.DB) error {
	// Проверяем, существует ли индекс
	var count int64
	db.Raw("SELECT COUNT(*) FROM INFORMATION_SCHEMA.STATISTICS WHERE table_schema = DATABASE() AND table_name = 'vacancy' AND index_name = 'idx_vacancy_fulltext'").Scan(&count)

	if count == 0 {
		// Создаем FULLTEXT индекс
		if err := db.Exec("ALTER TABLE vacancy ADD FULLTEXT INDEX idx_vacancy_fulltext (title, description)").Error; err != nil {
			return err
		}
		fmt.Println("FULLTEXT index created successfully")
	} else {
		fmt.Println("FULLTEXT index already exists")
	}

	return nil
}

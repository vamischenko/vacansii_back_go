package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// Config структура конфигурации приложения
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	CORS     CORSConfig
	RateLimit RateLimitConfig
}

// ServerConfig конфигурация сервера
type ServerConfig struct {
	Port string
	Mode string
}

// DatabaseConfig конфигурация базы данных
type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	Charset  string
}

// CORSConfig конфигурация CORS
type CORSConfig struct {
	AllowOrigins []string
}

// RateLimitConfig конфигурация Rate Limiting
type RateLimitConfig struct {
	Requests int
	Window   int // в секундах
}

// Load загружает конфигурацию из .env файла
func Load() *Config {
	// Загружаем .env файл
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	return &Config{
		Server: ServerConfig{
			Port: getEnv("SERVER_PORT", "8080"),
			Mode: getEnv("GIN_MODE", "debug"),
		},
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "3306"),
			User:     getEnv("DB_USER", "root"),
			Password: getEnv("DB_PASSWORD", ""),
			DBName:   getEnv("DB_NAME", "vakansii_db"),
			Charset:  getEnv("DB_CHARSET", "utf8mb4"),
		},
		CORS: CORSConfig{
			AllowOrigins: []string{getEnv("CORS_ORIGIN", "http://localhost:3000")},
		},
		RateLimit: RateLimitConfig{
			Requests: getEnvAsInt("RATE_LIMIT_REQUESTS", 100),
			Window:   getEnvAsInt("RATE_LIMIT_WINDOW", 3600),
		},
	}
}

// getEnv получает переменную окружения или возвращает значение по умолчанию
func getEnv(key string, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

// getEnvAsInt получает переменную окружения как integer или возвращает значение по умолчанию
func getEnvAsInt(key string, defaultValue int) int {
	valueStr := getEnv(key, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return defaultValue
}

// GetDSN возвращает строку подключения к базе данных
func (c *Config) GetDSN() string {
	return c.Database.User + ":" + c.Database.Password + "@tcp(" +
		c.Database.Host + ":" + c.Database.Port + ")/" +
		c.Database.DBName + "?charset=" + c.Database.Charset +
		"&parseTime=True&loc=Local"
}

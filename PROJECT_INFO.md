# Информация о проекте

## Описание

Backend API сервис для управления вакансиями, реализованный на Go с использованием фреймворка Gin. Проект является портом оригинального PHP/Yii2 приложения на Go с сохранением всех основных функций и архитектуры.

## Технический стек

### Backend
- **Язык:** Go 1.23
- **Web Framework:** Gin
- **ORM:** GORM
- **База данных:** MySQL 8.0
- **Конфигурация:** godotenv

### Middleware
- **CORS:** gin-contrib/cors
- **Rate Limiting:** ulule/limiter

### DevOps
- **Контейнеризация:** Docker, Docker Compose
- **CI/CD:** Готов к интеграции

## Архитектура

Проект следует многослойной архитектуре (Clean Architecture):

```
┌─────────────────────────────────────────┐
│           HTTP Handlers                  │
│         (Controllers Layer)              │
└──────────────┬──────────────────────────┘
               │
┌──────────────▼──────────────────────────┐
│         Business Logic                   │
│         (Services Layer)                 │
└──────────────┬──────────────────────────┘
               │
┌──────────────▼──────────────────────────┐
│       Data Access Layer                  │
│      (Repositories Layer)                │
└──────────────┬──────────────────────────┘
               │
┌──────────────▼──────────────────────────┐
│           Database                       │
│         (MySQL/GORM)                     │
└──────────────────────────────────────────┘
```

### Слои приложения

#### 1. Controllers (Контроллеры)
- Обработка HTTP запросов
- Валидация входных данных
- Формирование HTTP ответов
- Минимальная бизнес-логика

**Файлы:**
- `controllers/vacancy_controller.go`

#### 2. Services (Сервисы)
- Бизнес-логика приложения
- Валидация данных
- Координация между репозиториями
- Обработка ошибок

**Файлы:**
- `services/vacancy_service.go`

#### 3. Repositories (Репозитории)
- Работа с базой данных
- CRUD операции
- Запросы через GORM
- Абстракция доступа к данным

**Файлы:**
- `repositories/vacancy_repository.go`

#### 4. Models (Модели)
- Определение структур данных
- Валидация на уровне модели
- GORM теги для маппинга

**Файлы:**
- `models/vacancy.go`
- `models/user.go`

#### 5. Middleware
- CORS
- Rate Limiting
- Логирование (встроено в Gin)

**Файлы:**
- `middleware/rate_limiter.go`

## Основные возможности

### 1. CRUD операции
- ✅ Создание вакансий
- ✅ Получение списка вакансий
- ✅ Получение конкретной вакансии
- ✅ Обновление вакансий
- ✅ Удаление вакансий

### 2. Полнотекстовый поиск
- Использование FULLTEXT индексов MySQL
- Поиск по заголовку и описанию
- Ранжирование по релевантности
- Пагинация результатов

### 3. Пагинация
- 10 записей на страницу
- Информация о количестве страниц
- Общее количество записей

### 4. Сортировка
- По зарплате (asc/desc)
- По дате создания (asc/desc)
- По релевантности (для поиска)

### 5. Выборочные поля
- Возможность запросить только нужные поля
- Оптимизация трафика
- Гибкость API

### 6. Rate Limiting
- 100 запросов в час на IP
- Защита от злоупотреблений
- Заголовки с информацией о лимите

### 7. CORS
- Настраиваемые разрешенные origins
- Поддержка всех HTTP методов
- Готовность к работе с frontend

### 8. Дополнительные поля
- Хранение в JSON формате
- Гибкость структуры данных
- Валидация размера (до 5000 символов)

## Структура базы данных

### Таблица: vacancy

```sql
CREATE TABLE `vacancy` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `title` varchar(255) NOT NULL,
  `description` text NOT NULL,
  `salary` bigint NOT NULL,
  `additional_fields` json DEFAULT NULL,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  FULLTEXT KEY `idx_vacancy_fulltext` (`title`,`description`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
```

### Таблица: user

```sql
CREATE TABLE `user` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `username` varchar(255) NOT NULL,
  `email` varchar(255) NOT NULL,
  `password_hash` varchar(255) NOT NULL,
  `auth_key` varchar(32) DEFAULT NULL,
  `access_token` varchar(64) DEFAULT NULL,
  `status` bigint DEFAULT 10,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `username` (`username`),
  UNIQUE KEY `email` (`email`),
  UNIQUE KEY `access_token` (`access_token`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
```

## API Endpoints

### Vacancy API

```
GET    /vacancy          - Получить список вакансий
GET    /vacancy/:id      - Получить вакансию по ID
POST   /vacancy          - Создать новую вакансию
PUT    /vacancy/:id      - Обновить вакансию
DELETE /vacancy/:id      - Удалить вакансию
GET    /vacancy/search   - Полнотекстовый поиск
```

## Конфигурация

### Переменные окружения (.env)

```env
# Сервер
SERVER_PORT=8080          # Порт сервера
GIN_MODE=debug            # Режим: debug, release

# База данных
DB_HOST=localhost         # Хост MySQL
DB_PORT=3306              # Порт MySQL
DB_USER=root              # Пользователь
DB_PASSWORD=              # Пароль
DB_NAME=vakansii_db       # Имя БД
DB_CHARSET=utf8mb4        # Кодировка

# CORS
CORS_ORIGIN=http://localhost:3000  # Разрешенный origin

# Rate Limiting
RATE_LIMIT_REQUESTS=100   # Количество запросов
RATE_LIMIT_WINDOW=3600    # Окно в секундах (1 час)
```

## Миграции

Миграции выполняются автоматически при запуске приложения с использованием GORM AutoMigrate.

### Что создается:
1. Таблицы: `vacancy`, `user`
2. FULLTEXT индекс на поля `title` и `description`

## Развертывание

### Локальное развертывание

```bash
# 1. Установить зависимости
go mod download

# 2. Настроить .env
cp .env.example .env

# 3. Создать БД
mysql -u root -p
CREATE DATABASE vakansii_db;

# 4. Запустить
go run main.go
```

### Docker развертывание

```bash
# Запустить все сервисы
docker-compose up -d

# API: http://localhost:8080
# phpMyAdmin: http://localhost:8081
```

## Производительность

### Оптимизации

1. **FULLTEXT индексы** - быстрый поиск (10-100x быстрее LIKE)
2. **Пагинация** - ограничение нагрузки
3. **GORM** - эффективные SQL запросы
4. **Индексы БД** - быстрый доступ к данным

### Масштабирование

- Горизонтальное: несколько экземпляров приложения
- Вертикальное: увеличение ресурсов сервера
- База данных: репликация MySQL

## Безопасность

1. **Rate Limiting** - защита от злоупотреблений
2. **Валидация** - проверка входных данных
3. **Хеширование паролей** - bcrypt
4. **CORS** - контроль доступа
5. **SQL Injection** - защита через GORM

## Тестирование

```bash
# Запустить все тесты
go test -v ./...

# Тесты с покрытием
go test -v -cover ./...
```

## Сравнение с оригиналом (Yii2)

| Характеристика | Yii2 (PHP) | Gin (Go) |
|----------------|------------|----------|
| Производительность | Средняя | Высокая |
| Потребление памяти | Среднее | Низкое |
| Конкурентность | Ограничена | Отличная (goroutines) |
| Типобезопасность | Слабая | Сильная |
| Компиляция | Нет | Да |
| Размер бинарника | N/A | ~15 MB |
| Время старта | ~100ms | ~10ms |

## Дальнейшее развитие

### Планируемые улучшения

1. ✅ JWT аутентификация
2. ✅ Swagger документация
3. ✅ Unit тесты
4. ✅ Integration тесты
5. ✅ Кеширование (Redis)
6. ✅ Логирование (Zap)
7. ✅ Метрики (Prometheus)
8. ✅ CI/CD (GitHub Actions)

## Лицензия

MIT

## Контакты

Вопросы и предложения: [GitHub Issues]

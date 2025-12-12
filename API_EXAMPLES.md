# API Примеры использования

Примеры curl запросов для работы с API вакансий.

## Базовый URL

```
http://localhost:8080
```

## Получение списка вакансий

### Получить первую страницу (по умолчанию)

```bash
curl http://localhost:8080/vacancy
```

### Получить вторую страницу

```bash
curl "http://localhost:8080/vacancy?page=2"
```

### Сортировка по зарплате (по убыванию)

```bash
curl "http://localhost:8080/vacancy?sort=salary&order=desc"
```

### Сортировка по дате создания (по возрастанию)

```bash
curl "http://localhost:8080/vacancy?sort=created_at&order=asc"
```

## Получение конкретной вакансии

### Получить все поля вакансии с ID 1

```bash
curl http://localhost:8080/vacancy/1
```

### Получить только определенные поля

```bash
curl "http://localhost:8080/vacancy/1?fields=title,salary"
```

```bash
curl "http://localhost:8080/vacancy/1?fields=title,description,salary,additional_fields"
```

## Создание вакансии

### Создать простую вакансию

```bash
curl -X POST http://localhost:8080/vacancy \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Go Developer",
    "description": "Требуется опытный Go разработчик с знанием Gin",
    "salary": 180000
  }'
```

### Создать вакансию с дополнительными полями

```bash
curl -X POST http://localhost:8080/vacancy \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Senior Go Developer",
    "description": "Требуется опытный Senior Go разработчик",
    "salary": 250000,
    "additional_fields": {
      "company": "Tech Corp",
      "location": "Москва",
      "remote": true,
      "experience": "5+ лет",
      "skills": ["Go", "Gin", "MySQL", "Docker"]
    }
  }'
```

## Обновление вакансии

### Обновить зарплату

```bash
curl -X PUT http://localhost:8080/vacancy/1 \
  -H "Content-Type: application/json" \
  -d '{
    "salary": 200000
  }'
```

### Обновить несколько полей

```bash
curl -X PUT http://localhost:8080/vacancy/1 \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Senior Go Developer",
    "salary": 220000,
    "description": "Обновленное описание вакансии"
  }'
```

### Обновить дополнительные поля

```bash
curl -X PUT http://localhost:8080/vacancy/1 \
  -H "Content-Type: application/json" \
  -d '{
    "additional_fields": {
      "company": "New Tech Corp",
      "location": "Санкт-Петербург",
      "remote": false
    }
  }'
```

## Удаление вакансии

```bash
curl -X DELETE http://localhost:8080/vacancy/1
```

## Полнотекстовый поиск

### Поиск по одному слову

```bash
curl "http://localhost:8080/vacancy/search?q=Go"
```

### Поиск по нескольким словам

```bash
curl "http://localhost:8080/vacancy/search?q=Go%20разработчик"
```

### Поиск с сортировкой по релевантности (по умолчанию)

```bash
curl "http://localhost:8080/vacancy/search?q=developer&sort=relevance"
```

### Поиск с сортировкой по дате

```bash
curl "http://localhost:8080/vacancy/search?q=PHP&sort=desc"
```

### Поиск с пагинацией

```bash
curl "http://localhost:8080/vacancy/search?q=developer&page=2"
```

## Примеры ответов

### Успешное получение списка

```json
{
  "data": [
    {
      "id": 1,
      "title": "Go Developer",
      "description": "Требуется опытный Go разработчик",
      "salary": 180000,
      "additional_fields": {
        "company": "Tech Corp",
        "location": "Москва"
      },
      "created_at": "2025-01-10T12:00:00Z",
      "updated_at": "2025-01-10T12:00:00Z"
    }
  ],
  "pagination": {
    "total": 25,
    "page": 1,
    "pageSize": 10,
    "pageCount": 3
  }
}
```

### Успешное создание

```json
{
  "success": true,
  "id": 13,
  "message": "Вакансия успешно создана"
}
```

### Ошибка валидации

```json
{
  "success": false,
  "message": "Название вакансии обязательно"
}
```

### Вакансия не найдена

```json
{
  "success": false,
  "message": "Вакансия не найдена"
}
```

## Rate Limiting

API ограничивает количество запросов до 100 в час с одного IP адреса.

### Превышение лимита

```json
{
  "success": false,
  "message": "Превышен лимит запросов. Попробуйте позже."
}
```

## Скрипт для тестирования

```bash
#!/bin/bash

BASE_URL="http://localhost:8080"

echo "=== Создание вакансии ==="
VACANCY_ID=$(curl -s -X POST $BASE_URL/vacancy \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Test Vacancy",
    "description": "Test Description",
    "salary": 100000
  }' | jq -r '.id')

echo "Создана вакансия с ID: $VACANCY_ID"

echo -e "\n=== Получение вакансии ==="
curl -s "$BASE_URL/vacancy/$VACANCY_ID" | jq

echo -e "\n=== Обновление вакансии ==="
curl -s -X PUT "$BASE_URL/vacancy/$VACANCY_ID" \
  -H "Content-Type: application/json" \
  -d '{"salary": 150000}' | jq

echo -e "\n=== Поиск вакансий ==="
curl -s "$BASE_URL/vacancy/search?q=Test" | jq

echo -e "\n=== Удаление вакансии ==="
curl -s -X DELETE "$BASE_URL/vacancy/$VACANCY_ID" | jq
```

Сохраните этот скрипт как `test_api.sh` и запустите:

```bash
chmod +x test_api.sh
./test_api.sh
```

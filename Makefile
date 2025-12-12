.PHONY: help build run test clean docker-up docker-down docker-logs

help:
	@echo "Доступные команды:"
	@echo "  make build        - Собрать приложение"
	@echo "  make run          - Запустить приложение"
	@echo "  make test         - Запустить тесты"
	@echo "  make clean        - Очистить собранные файлы"
	@echo "  make docker-up    - Запустить Docker контейнеры"
	@echo "  make docker-down  - Остановить Docker контейнеры"
	@echo "  make docker-logs  - Показать логи Docker контейнеров"

build:
	@echo "Сборка приложения..."
	go build -o main .

run:
	@echo "Запуск приложения..."
	go run main.go

test:
	@echo "Запуск тестов..."
	go test -v ./...

clean:
	@echo "Очистка..."
	rm -f main
	go clean

docker-up:
	@echo "Запуск Docker контейнеров..."
	docker-compose up -d

docker-down:
	@echo "Остановка Docker контейнеров..."
	docker-compose down

docker-logs:
	@echo "Логи Docker контейнеров..."
	docker-compose logs -f

tidy:
	@echo "Обновление зависимостей..."
	go mod tidy

fmt:
	@echo "Форматирование кода..."
	go fmt ./...

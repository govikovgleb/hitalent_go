.PHONY: help up seed down build logs test migrate

# Default goal
.DEFAULT_GOAL := help

# Colors for output
BLUE := \033[36m
GREEN := \033[32m
YELLOW := \033[33m
RED := \033[31m
NC := \033[0m # No Color

help: ## Показать список команд
	@echo "$(BLUE)Доступные команды:$(NC)"
	@echo ""
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  $(GREEN)%-15s$(NC) %s\n", $$1, $$2}'
	@echo ""

up: ## Запустить все сервисы (app + db + migrations)
	@echo "$(BLUE)Запуск приложения...$(NC)"
	docker-compose up -d --build
	@echo "$(GREEN)Приложение запущено на http://localhost:8080$(NC)"

seed: ## Заполнить базу тестовыми данными (отдельная команда)
	@echo "$(BLUE)Заполнение базы тестовыми данными...$(NC)"
	@docker-compose run --rm -e DATABASE_URL="postgres://talent:talentpass@db:5432/hitalent?sslmode=disable" app go run ./cmd/seed
	@echo "$(GREEN)Тестовые данные добавлены!$(NC)"
	@echo "$(YELLOW)Проверь структуру: curl http://localhost:8080/departments/1?depth=5$(NC)"
	@echo "$(YELLOW)Удали IT каскадно: curl -X DELETE 'http://localhost:8080/departments/2?mode=cascade'$(NC)"

down: ## Остановить все сервисы
	@echo "$(BLUE)Остановка сервисов...$(NC)"
	docker-compose down

build: ## Пересобрать контейнеры
	@echo "$(BLUE)Пересборка контейнеров...$(NC)"
	docker-compose build --no-cache

logs: ## Показать логи приложения
	docker-compose logs -f app

logs-all: ## Показать логи всех сервисов
	docker-compose logs -f

test: ## Запустить тесты
	@echo "$(BLUE)Запуск тестов...$(NC)"
	go test ./internal/service/ -v

migrate: ## Выполнить миграции БД
	@echo "$(BLUE)Выполнение миграций...$(NC)"
	docker-compose run --rm migrations goose up


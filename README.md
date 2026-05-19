# API Организационной структуры

REST API для управления организационной структурой компании с поддержкой иерархических подразделений и сотрудников.

## Стек технологий

- **Go 1.23** - язык программирования
- **GORM** - ORM для работы с PostgreSQL
- **PostgreSQL** - база данных
- **Goose** - инструмент миграций
- **Docker & Docker Compose** - контейнеризация
- **go-playground/validator** - валидация данных

## Быстрый старт

### Запуск

```bash
# Клонирование и запуск
git clone https://github.com/govikovgleb/hitalent_go.git
cd hitalent_go
make up
```

Приложение будет доступно на `http://localhost:8080`

### Запуск с тестовыми данными (Seeder)

Для удобного тестирования API есть **отдельная команда** `make seed` которая заполняет базу тестовыми данными:

```bash
# 1. Сначала запускаем базу и миграции
make up

# 2. Заполняем тестовыми данными
make seed
```

Seeder работает как отдельная утилита (`cmd/seed/main.go`).

**Создаваемая структура:**
```
Компания (ID=1)
├── IT (ID=2) - 3 сотрудника
│   ├── Backend (ID=3) - 4 сотрудника
│   │   └── Go Team (ID=4) - 2 сотрудника
│   ├── Frontend (ID=5) - 3 сотрудника
│   └── DevOps (ID=6) - 2 сотрудника
├── HR (ID=7) - 2 сотрудника
└── Sales (ID=8) - 2 сотрудника
    ├── B2B Sales (ID=9) - 3 сотрудника
    └── B2C Sales (ID=10) - 3 сотрудника

Всего: 10 департаментов, 26 сотрудников
```

**Примеры тестирования после seeding:**

```bash
# Посмотреть всю структуру
curl "http://localhost:8080/departments/1?depth=5&include_employees=true" | jq .

# Удалить IT департамент каскадно
curl -X DELETE "http://localhost:8080/departments/2?mode=cascade"

# Удалить Sales и перевести сотрудников в HR
curl -X DELETE "http://localhost:8080/departments/8?mode=reassign&reassign_dept_id=7"
```

### Остановка

```bash
make down
```

## API Endpoints

### Departments

#### Создать подразделение
```bash
POST /departments/
Content-Type: application/json

{
  "name": "Engineering",
  "parent_id": null  // или ID родительского подразделения
}
```

#### Получить подразделение
```bash
GET /departments/{id}?depth=2&include_employees=true
```

Параметры:
- `depth` (1-5) - глубина вложенных подразделений, по умолчанию 1
- `include_employees` (true/false) - включать сотрудников, по умолчанию true

#### Обновить подразделение
```bash
PATCH /departments/{id}
Content-Type: application/json

{
  "name": "New Name",
  "parent_id": 5  // или null для корневого
}
```

#### Удалить подразделение
```bash
DELETE /departments/{id}?mode=cascade
DELETE /departments/{id}?mode=reassign&reassign_dept_id=5
```

Режимы:
- `cascade` - удалить подразделение, всех сотрудников и дочерние подразделения
- `reassign` - удалить подразделение, сотрудников перевести в указанное

### Сотрудники (Employees)

#### Создать сотрудника
```bash
POST /departments/{id}/employees/
Content-Type: application/json

{
  "full_name": "Иванов Иван Иванович",
  "position": "Senior Developer",
  "hired_at": "2026-05-15"  // опционально
}
```

## Примеры использования

### Создание структуры компании

```bash
# Создаем корневое подразделение
curl -X POST http://localhost:8080/departments/ \
  -H "Content-Type: application/json" \
  -d '{"name": "Компания", "parent_id": null}'

# Создаем IT отдел
curl -X POST http://localhost:8080/departments/ \
  -H "Content-Type: application/json" \
  -d '{"name": "IT", "parent_id": 1}'

# Добавляем сотрудника в IT отдел
curl -X POST http://localhost:8080/departments/2/employees/ \
  -H "Content-Type: application/json" \
  -d '{"full_name": "Петров Петр", "position": "Разработчик"}'

# Получаем структуру с подразделениями и сотрудниками
curl "http://localhost:8080/departments/1?depth=3&include_employees=true"
```

## Разработка

### Структура проекта

```
.
├── cmd/
│   ├── api/                # Точка входа приложения
│   └── seed/               # Утилита для заполнения тестовыми данными
├── internal/
│   ├── config/             # Конфигурация
│   ├── handlers/           # HTTP handlers
│   ├── models/             # Модели данных (GORM)
│   ├── repository/         # Работа с БД
│   ├── service/            # Бизнес-логика
│   ├── validator/          # Валидация
│   ├── seed/               # Seeder для тестовых данных
│   └── router/             # Маршрутизация
├── migrations/             # SQL миграции (goose)
├── docker/
│   └── Dockerfile          # Многостадийная сборка
├── docker-compose.yml      # Docker Compose конфиг
├── Makefile               # Команды для разработки
└── README.md
```

### Database Seeding

Для быстрого наполнения базы тестовыми данными используется **seeder** (internal/seed/seed.go).

#### Что создает seeder:
- **Иерархия департаментов**: 3 уровня вложенности (Компания → IT → Backend → Go Team)
- **Сотрудники**: 26 сотрудников в разных отделах с разными должностями
- **ID**: Департаменты получают ID 1-10 (последовательно при первом запуске)

#### Режимы запуска:
**Через Makefile** : `make seed`

### Команды Makefile

```bash
make help       # Показать все команды
make up         # Запустить все сервисы
make down       # Остановить все сервисы
make build      # Пересобрать контейнеры
make logs       # Показать логи
make test       # Запустить тесты
make seed       # Заполнить тестовыми данными
make migrate    # Выполнить миграции
make migrate-down # Откатить миграции
make clean      # Очистить все (контейнеры + volumes)
```

### Запуск тестов

```bash
make test
```

## Архитектура

Проект построен с разделением на слои:

- **Handler** - HTTP обработчики, сериализация JSON
- **Service** - бизнес-логика, работа с деревом
- **Repository** - работа с БД через GORM
- **Models** - определение сущностей и связей
- **Validator** - валидация

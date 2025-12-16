# Gym StrongCode - Fitness Club Management System

Backend REST API для управления фитнес-клубом на Go.

## Особенности

- JWT аутентификация + роли (user/admin)
- Полный CRUD для всех сущностей (Users, Gyms, Trainers, Classes, Memberships, Bookings, Payments)
- Отношения в БД: one-to-many (Trainer → Classes, Gym → Classes, User → Bookings/Payments), many-to-many (User ↔ Memberships)
- Уведомления по email (Gmail SMTP, асинхронно через background worker)
- Structured logging в файл `logs/app.log` (JSON-формат, легко искать)
- Graceful shutdown и propagation context
- Rate limiting middleware
- Request logging middleware
- Swagger API документация
- Автоматические миграции (golang-migrate) с seed-данными
- Тесты (unit + integration)
- Docker + docker-compose с persistent volumes (база и логи сохраняются)

## Архитектура

- Gin – HTTP фреймворк
- SQLite – легковесная БД (с foreign keys и индексами)
- golang-migrate – миграции
- JWT (golang-jwt) – аутентификация
- Zap – структурированное логирование
- Viper – загрузка конфигурации (.env)
- Swaggo – Swagger docs
- Background goroutines – очередь отправки email

## Требования

- Go 1.23+
- Docker & Docker Compose (рекомендуется)
- Gmail аккаунт (для SMTP уведомлений – используйте App Password)

## Установка и запуск

### Вариант 1: Docker (рекомендуется)

```bash
# Клонируем репозиторий
git clone https://github.com/ZhigerDinmukhamed/Gym_StrongCode.git
cd Gym_StrongCode

# Копируем и настраиваем .env
cp .env.example .env

# Опционально: настройте SMTP в .env для email-уведомлений
# SMTP_HOST=smtp.gmail.com
# SMTP_PORT=587
# SMTP_USER=your@gmail.com
# SMTP_PASS=your-app-password
# FROM_EMAIL=your@gmail.com

# Запускаем
docker-compose up -d --build

# API доступен по http://localhost:8080
# Swagger UI: http://localhost:8080/swagger/index.html
# База сохраняется в ./data/gym_strongcode.db
# Логи сохраняются в ./logs/app.log
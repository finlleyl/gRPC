# Сервис SSO на gRPC

Минималистичный **Single Sign‑On (SSO)** микросервис на Go, предоставляющий gRPC‑API для **регистрации пользователей** и **аутентификации** с выдачей подписанных **JWT**‑токенов.

<p align="center">
  <img src="https://user-images.githubusercontent.com/20046/273320805-a6d89aa8-e2c4-4024-9bce-c26d15ff7eed.png" width="520" alt="Диаграмма архитектуры" />
</p>

---

## 📜 Содержание

1. [Возможности](#-возможности)
2. [Стек технологий](#-стек-технологий)
3. [Структура проекта](#-структура-проекта)
4. [Быстрый старт](#-быстрый-старт)

   * [Зависимости](#зависимости)
   * [Генерация Protobuf](#генерация-protobuf)
   * [Запуск миграций](#запуск-миграций)
   * [Запуск сервиса](#запуск-сервиса)
5. [Конфигурация](#-конфигурация)
6. [gRPC API](#-grpc-api)
7. [Пример использования](#-пример-использования)
8. [Формат токена](#-формат-токена)
9. [Схема базы данных](#-схема-базы-данных)
10. [Вклад](#-вклад)
11. [Лицензия](#-лицензия)

---

## ✨ Возможности

| Возможность          | Описание                                                       |
| -------------------- | -------------------------------------------------------------- |
| **Регистрация**      | Создание пользователя по *email* с хешированным паролем        |
| **Вход**             | Проверка учётных данных, выдача JWT для заданного **app\_id**  |
| **SQLite**‑хранилище | Лёгкая файл‑БД для пользователей и приложений                  |
| **Миграции**         | Управление схемой через `golang-migrate`, CLI в `cmd/migrator` |
| **Генерация JWT**    | HS256, настраиваемый TTL, стандартные + кастомные claims       |
| **Структурные логи** | zap‑логгер (dev/prod режимы)                                   |
| **Taskfile**         | Единственная команда для генерации stубов и рутинных задач     |

---

## 🔧 Стек технологий

* **Go 1.24**
* **gRPC 1.60** и **Protocol Buffers v3**
* **SQLite 3** (`github.com/mattn/go-sqlite3`)
* **golang-migrate** для миграций
* **zap** для логирования
* **bcrypt** для хеширования паролей

---

## 🗂 Структура проекта

```
.
├── cmd/            # Точки входа (migrator и gRPC‑сервер)
├── config/         # YAML‑конфиги
├── internal/       # Домены, сервисы, транспорты
│   ├── app/        # DI и инициализация
│   ├── grpc/       # gRPC‑обёртка
│   ├── services/   # Бизнес‑логика (auth)
│   ├── storage/    # Репозитории (SQLite)
│   └── lib/        # Пакеты общего назначения (jwt и др.)
├── proto/          # .proto схемы
├── gen/            # Сгенерированные стабы (не коммитятся)
├── migrations/     # SQL‑миграции
└── Taskfile.yaml   # Task‑раннер
```

---

## 🚀 Быстрый старт

### Зависимости

| Инструмент                | Версия  | Как установить                                                    |
| ------------------------- | ------- | ----------------------------------------------------------------- |
| **Go**                    | ≥ 1.24  | [https://go.dev/dl/](https://go.dev/dl/)                          |
| **protoc**                | ≥ 3.21  | Пакетный менеджер или release‑архив                               |
| **protoc-gen-go**         | latest  | `go install google.golang.org/protobuf/cmd/protoc-gen-go@latest`  |
| **protoc-gen-go-grpc**    | latest  | `go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest` |
| **Task** (опц.)           | ≥ 3     | `go install github.com/go-task/task/v3/cmd/task@latest`           |
| **golang-migrate** (опц.) | ≥ v4.16 | нужен, если не использовать `cmd/migrator`                        |

> **Подсказка:** `task check-tools` проверит и при необходимости установит плагины для Protobuf.

---

### Генерация Protobuf

```bash
# Скрипт (рекурсивно ищет .proto и генерирует stubs в gen/go)
$ task generate
```

Ручной вызов:

```bash
protoc -I proto proto/sso/sso.proto \
  --go_out=gen/go --go_opt=paths=source_relative \
  --go-grpc_out=gen/go --go-grpc_opt=paths=source_relative
```

---

### Запуск миграций

```bash
# Сборка утилиты
$ go build -o bin/migrator ./cmd/migrator

# Применить миграции (БД создаётся при первом запуске)
$ ./bin/migrator \
  --storage-path=./storage/sso.db \
  --migrations-path=./migrations \
  --migrations-table=migrations
```

> Пути по умолчанию соответствуют *config/config\_local.yaml*.

---

### Запуск сервиса

```bash
# Сборка бинаря сервера
$ go build -o bin/sso ./cmd/sso

# Запуск с явным конфигом
$ ./bin/sso -config=config/config_local.yaml

# ИЛИ через переменную окружения
$ CONFIG_PATH=config/config_local.yaml ./bin/sso
```

В логах появится:

```
INFO	starting gRPC server on port	{"port":44044}
```

---

## ⚙️ Конфигурация

Параметры берутся из YAML‑файла; путь задаётся флагом `-config` или переменной `CONFIG_PATH`.

```yaml
# config/config_local.yaml
env: "local"
storage_path: "./storage/sso.db"
grpc:
  port: 44044      # Порт прослушивания
  timeout: 10h     # Дедлайн на запрос

token_ttl: 1h      # Время жизни JWT
```

Для prod‑окружения можно подготовить отдельный файл (например `config/config_prod.yaml`).

---

## 🛰️ gRPC API

Определение ─ [`proto/sso/sso.proto`](proto/sso/sso.proto):

```protobuf
service Auth {
  rpc Register (RegisterRequest) returns (RegisterResponse);
  rpc Login    (LoginRequest)    returns (LoginResponse);
}
```

| Метод      | Параметры запроса             | Ответ         |
| ---------- | ----------------------------- | ------------- |
| `Register` | `email`, `password`           | `user_id`     |
| `Login`    | `email`, `password`, `app_id` | `token` (JWT) |

> **Коды ошибок** соответствуют gRPC‑стандарту. Важные случаи:
>
> * `InvalidArgument` – отсутствуют/некорректны поля
> * `AlreadyExists`  – пользователь уже зарегистрирован
> * `Internal`       – внутренняя ошибка сервера

---

## 🔌 Пример использования

### grpcurl

```bash
# Регистрация
grpcurl -plaintext -d '{"email":"bob@example.com","password":"s3cr3t"}' \
  localhost:44044 auth.Auth/Register

# Логин
grpcurl -plaintext -d '{"email":"bob@example.com","password":"s3cr3t","app_id":1}' \
  localhost:44044 auth.Auth/Login
```

### Go‑клиент

```go
conn, _ := grpc.Dial("localhost:44044", grpc.WithInsecure())
client := ssopb.NewAuthClient(conn)

resp, err := client.Login(ctx, &ssopb.LoginRequest{
    Email:    "bob@example.com",
    Password: "s3cr3t",
    AppId:    1,
})
fmt.Println("JWT:", resp.Token)
```

---

## 🔐 Формат токена

JWT подписываются алгоритмом **HS256** с секретом, уникальным для приложения (таблица `apps`). Основные claim‑ы:

| Claim    | Тип    | Описание                              |
| -------- | ------ | ------------------------------------- |
| `uid`    | int64  | ID пользователя                       |
| `email`  | string | Email пользователя                    |
| `app_id` | int64  | ID приложения, запросившего токен     |
| `exp`    | int64  | Время истечения (`now` + `token_ttl`) |

---

## 🗄️ Схема базы данных

```sql
CREATE TABLE IF NOT EXISTS users (
    id         INTEGER PRIMARY KEY,
    email      TEXT NOT NULL UNIQUE,
    pass_hash  BLOB NOT NULL
);

CREATE TABLE IF NOT EXISTS apps (
    id     INTEGER PRIMARY KEY,
    name   TEXT NOT NULL UNIQUE,
    secret TEXT NOT NULL UNIQUE
);
```

Пример предзаполнения таблицы приложений:

```sql
INSERT INTO apps(name, secret) VALUES ('my-frontend', 'super-secret-key');
```

---
яется под лицензией **MIT**. Подробности см. в файле `LICENSE`.

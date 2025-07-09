# Go Animals

A sample RESTful API built with [Gin](https://github.com/gin-gonic/gin), [PostgreSQL](https://www.postgresql.org/) using [sqlx](https://github.com/jmoiron/sqlx).

This service manages an `animals` table in a `sampledb` PostgreSQL database.

---

## 🧰 Prerequisites

- Go 1.21+
- Docker

---

## 🚀 Running the application

1. Set up the environment:

```bash
make environment
cp .env.sample .env
```

2. Start the server:

```bash
make server
```

> The application will be available at `http://localhost:8080`

3. See all available commands:

```bash
make help
```

---

## 🐘 Expected table structure

```sql
sampledb=# SELECT * FROM animals;
 id | name | age | description
----+------+-----+-------------
(0 rows)
```

---

## 🎯 Example requests with `curl`

### 1. Create an animal

```bash
curl -X POST http://localhost:8080/animals \
  -H "Content-Type: application/json" \
  -d '{"name": "cow", "age": 20, "description": "beautiful cow"}'
```

### 2. List all animals

```bash
curl http://localhost:8080/animals
```

### 3. Get animal by ID

```bash
curl http://localhost:8080/animals/11
```

### 4. Update animal

```bash
curl -X PATCH http://localhost:8080/animals/12 \
  -H "Content-Type: application/json" \
  -d '{"name": "cat", "age": 15, "description": "beautiful cat update"}'
```

### 5. Delete animal

```bash
curl -X DELETE http://localhost:8080/animals/13
```

---

## ✅ Best Practices

### 🧱 Architecture

- Modular project layout (`internal/`, `router/`, `config/`)
- Clear separation of handler, repository, and domain logic
- Contract-first API development with OpenAPI

### 🪵 Structured Logging with Zap

- Single instance of logger injected via context or constructor
- Environment-based format: JSON in production, colorized console in development
- All success operations are logged using `logger.Info(...)`

### ⚠️ Error Handling

- All errors wrapped with `fmt.Errorf(..., %w)` for context
- Domain-level error types (e.g., `ErrAnimalNotFound`)
- Use of `errors.Is()` for reliable error type matching
- No panics used for expected application flow

### 🧪 Testing

- Unit tests with mocks using `testify/assert`
- End-to-end (E2E) tests using `testcontainers-go` + real PostgreSQL
- Tests verify real DB state after operations
- Code coverage via: `go test -coverprofile=coverage.out ./...`

### 🧬 Type Mapping

- PostgreSQL `BIGSERIAL` → Go `int64`

---

## 📚 References

- [Gin](https://github.com/gin-gonic/gin)
- [Zap Logger](https://github.com/uber-go/zap)
- [sqlx](https://github.com/jmoiron/sqlx)
- [Testcontainers-Go](https://github.com/testcontainers/testcontainers-go)

---

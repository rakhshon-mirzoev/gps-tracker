# Откуда берется `request_id`?

## 🔄 Полный поток

```
1️⃣ HTTP Запрос приходит
   ├─ Заголовок: X-Request-ID: "123e4567-e89b-12d3-a456-426614174000"  (опционально)
   └─ Или: генерируется новый UUID если заголовка нет

2️⃣ Middleware: RequestID()
   ├─ Парсит X-Request-ID из заголовка (или генерирует новый)
   ├─ Сохраняет в контекст: c.Set(constants.RequestIDKey, reqID)
   └─ Добавляет в ответ: c.Writer.Header().Set("X-Request-ID", reqID)

3️⃣ Handler выполняется
   └─ Может получить: c.GetString(constants.RequestIDKey)

4️⃣ Логирование в handler
   ├─ requestID := c.GetString(constants.RequestIDKey)  ← вот отсюда!
   └─ logger.Log.Error(..., zap.String("request_id", requestID))

5️⃣ HTTP Ответ
   ├─ JSON: {"request_id": "123e4567-e89b-12d3-a456-426614174000"}
   └─ Заголовок: X-Request-ID: "123e4567-e89b-12d3-a456-426614174000"
```

## 🔑 Ключевые компоненты

### 1. Middleware (internal/middleware/request_id.go)
```go
func (s *MiddlewareService) RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Получаем ID из заголовка ИЛИ генерируем новый
		reqID := c.GetHeader("X-Request-ID")
		if reqID == "" {
			reqID = uuid.New().String()
		}

		// ✅ Сохраняем в контекст
		c.Set(constants.RequestIDKey, reqID)
		
		// Добавляем в ответ
		c.Writer.Header().Set("X-Request-ID", reqID)

		c.Next()  // Передаем управление дальше
	}
}
```

### 2. Константа (pkg/constants/constants.go)
```go
const RequestIDKey = "request_id"  // Ключ для c.Set() и c.Get()
```

### 3. Handler получает (internal/handler/error_tracing.go)
```go
requestID := c.GetString(constants.RequestIDKey)  // Достаем из контекста
```

## 📍 Где это используется?

✅ **В приложениях уже подключено:**
- `/cmd/mobile/routes.go` - строка 16: `mw.RequestID()`
- `/cmd/admin/routes.go` - строка 16: `mw.RequestID()`

```go
// Пример из routes.go
func setupRoutes(engine *gin.Engine, mw middleware.Middleware) {
    engine.Use(
        mw.RequestID(),      // ← Это устанавливает request_id в контекст
        mw.Logger(),
        mw.CorsMiddleWare(),
    )
}
```

## 🎯 Как это работает шаг за шагом

### Запрос 1: С заголовком X-Request-ID
```
Запрос:
  GET /api/users/123
  X-Request-ID: my-custom-id-456

Middleware RequestID():
  reqID := c.GetHeader("X-Request-ID")  // "my-custom-id-456"
  c.Set("request_id", "my-custom-id-456")

Handler:
  requestID := c.GetString("request_id")  // "my-custom-id-456"

Логирование:
  {"request_id": "my-custom-id-456", ...}

Ответ:
  X-Request-ID: my-custom-id-456
```

### Запрос 2: Без заголовка X-Request-ID
```
Запрос:
  GET /api/users/123
  (нет X-Request-ID)

Middleware RequestID():
  reqID := c.GetHeader("X-Request-ID")  // пусто
  if reqID == "" {
    reqID = uuid.New().String()  // "a1b2c3d4-e5f6-..."
  }
  c.Set("request_id", "a1b2c3d4-e5f6-...")

Handler:
  requestID := c.GetString("request_id")  // "a1b2c3d4-e5f6-..."

Логирование:
  {"request_id": "a1b2c3d4-e5f6-...", ...}

Ответ:
  X-Request-ID: a1b2c3d4-e5f6-...
```

## 💡 Практическая пример отслеживания

Клиент делает запрос и хочет понять, что произошло:

```bash
# Запрос с явным ID
curl -H "X-Request-ID: user-trace-001" \
     http://localhost:8000/api/staff-users/deactivate

# Или без - получит UUID
curl http://localhost:8000/api/staff-users/deactivate
# Ответ: X-Request-ID: 550e8400-e29b-41d4-a716-446655440000
```

**В логах приложения (все логи этого запроса):**
```json
[
  {
    "timestamp": "2026-02-04T10:30:45Z",
    "operation": "DeactivateStaffUser",
    "request_id": "550e8400-e29b-41d4-a716-446655440000",
    "message": "Checking system user",
    "level": "info"
  },
  {
    "timestamp": "2026-02-04T10:30:45Z",
    "operation": "DeactivateStaffUser",
    "request_id": "550e8400-e29b-41d4-a716-446655440000",
    "message": "Parsing UUID from params",
    "level": "info"
  },
  {
    "timestamp": "2026-02-04T10:30:45Z",
    "operation": "DeactivateStaffUser",
    "request_id": "550e8400-e29b-41d4-a716-446655440000",
    "message": "DeactivateStaffUser: operation succeeded",
    "level": "info"
  }
]
```

**Все логи связаны одним request_id!** Легко отследить что произошло в одном запросе.

## ⚙️ Что уже на месте:

✅ `RequestID()` middleware - реализован и используется
✅ `constants.RequestIDKey` - определена константа
✅ `handleError()` и `handleFieldError()` - обновлены на использование константы

## 🚀 Что делать дальше:

1. ✅ Убедитесь что middleware используется в routes.go (уже есть)
2. ✅ Используйте `h.handleError()` в handler'ах (вместо ручного логирования)
3. ✅ Request ID будет автоматически в логах и ответах

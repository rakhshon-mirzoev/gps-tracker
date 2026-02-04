# Error Tracing в Handler'ах - Лучшие Практики

## Архитектура трассировки ошибок

```
Repository Layer (БД операции)
    ↓ fmt.Errorf("operation: %w", err)
    
Service Layer (бизнес-логика)
    ↓ может добавить контекст
    
Handler Layer (HTTP граница)
    ↓ Структурированное логирование + HTTP ответ
    
Client
```

## Слои и их ответственность

### 1. Repository Layer
**Добавляет контекст БД операции**
```go
// Хорошо ✓
err := r.DB.Pool.QueryRow(ctx, sql, id).Scan(&user)
if err != nil {
    return fmt.Errorf("get user by id: %w", err)  // Контекст операции
}

// Плохо ✗
return err  // Потеряется контекст
```

### 2. Service Layer (опционально)
**Может добавить контекст бизнес-логики**
```go
func (s *staffUserService) Deactivate(ctx context.Context, id uuid.UUID, dto model.DeactivateStaffUser) error {
    // Валидация
    if err := s.validateDeactivation(dto); err != nil {
        return fmt.Errorf("validate deactivation: %w", err)
    }
    
    // Вызов repo с обработкой
    if err := s.repo.Deactivate(ctx, id, dto); err != nil {
        // Ошибка уже обогащена в repository!
        // Просто пробросить выше, если нужен доп контекст:
        return fmt.Errorf("deactivate staff user in service: %w", err)
    }
    
    return nil
}
```

### 3. Handler Layer
**Структурированное логирование + HTTP ответ**

#### Вариант 1: С helper функцией (рекомендуется)
```go
func (h *Handler) DeactivateStaffUser(c *gin.Context) {
    var deactivate model.DeactivateStaffUser
    
    uid, err := uuid.Parse(c.Param("id"))
    if err != nil {
        h.handleFieldError(c, "DeactivateStaffUser", err)  // Валидация
        return
    }
    
    if err := h.s.StaffUser.Deactivate(c.Request.Context(), uid, deactivate); err != nil {
        h.handleError(c, "DeactivateStaffUser", err, http.StatusInternalServerError)
        return
    }
    
    c.JSON(http.StatusOK, gin.H{"message": "success"})
}
```

#### Вариант 2: Встроенная обработка (если нет helper)
```go
func (h *Handler) DeactivateStaffUser(c *gin.Context) {
    // ... валидация ...
    
    if err := h.s.StaffUser.Deactivate(c.Request.Context(), uid, deactivate); err != nil {
        // Оборачиваем ошибку с контекстом handler'а
        fullErr := fmt.Errorf(
            "handler.DeactivateStaffUser: user_id=%s, staff_id=%s: %w",
            user_id, uid, err,
        )
        
        logger.Log.Error(
            "Failed to deactivate staff user",
            zap.Error(fullErr),
            zap.String("request_id", c.GetString("request_id")),
        )
        
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": "failed to deactivate user",
        })
        return
    }
}
```

## Трассировка через слои

### Пример полного потока ошибки:

```
1. Repository возвращает:
   "get user by id: ERROR: no rows in result set"

2. Service может добавить:
   "deactivate staff user in service: get user by id: ERROR: no rows in result set"

3. Handler оборачивает для логирования:
   "handler.DeactivateStaffUser: user_id=uuid1, staff_id=uuid2: 
    deactivate staff user in service: get user by id: ERROR: no rows in result set"

4. Логирование с контекстом:
   {
       "message": "Failed to deactivate staff user",
       "error": "...полная цепь ошибок...",
       "request_id": "req-123",
       "operation": "DeactivateStaffUser",
       "user_id": "uuid1",
       "status_code": 404
   }
```

## Правила в Handler'е

✓ **ПРАВИЛЬНО:**
- Использовать `handleError()` для бизнес-ошибок (NotFound, Validation)
- Использовать `handleFieldError()` для JSON/UUID парсинга
- Логировать с контекстом: operation, user_id, request_id
- Отправлять клиенту только безопасные сообщения (не стек трассы)

✗ **НЕПРАВИЛЬНО:**
- `c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})` - полная ошибка видна клиенту
- Логировать без контекста операции
- Игнорировать ошибки
- Множественные `if err != nil` без унификации

## Структура логируемой информации

```go
logger.Log.Error(
    "Operation description",           // Человеко-читаемое описание
    zap.Error(err),                    // Оригинальная ошибка (полный стек!)
    zap.String("operation", "OperationName"),  // Название операции
    zap.String("request_id", requestID),       // Трассировка запроса
    zap.String("user_id", userID),             // Контекст пользователя
    zap.String("resource_id", resourceID),     // Контекст ресурса
    zap.Int("status_code", statusCode),        // HTTP статус
)
```

## Request ID для трассировки

Добавьте middleware в app initialization:

```go
func (a *App) setupMiddleware() {
    a.engine.Use(func(c *gin.Context) {
        requestID := c.GetString("X-Request-ID")
        if requestID == "" {
            requestID = uuid.New().String()
        }
        c.Set("request_id", requestID)
        c.Next()
    })
}
```

Так все логи будут связаны одним request ID!

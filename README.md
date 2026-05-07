# GPS Tracker — API Documentation

Система мониторинга полевых сотрудников с GPS-трекингом, управлением визитами и планированием маршрутов.

---

## Стек

| Компонент | Технология |
|---|---|
| Framework | Gin (Go) |
| База данных | PostgreSQL |
| Хранилище файлов | MinIO |
| Очередь сообщений | Kafka |
| Push-уведомления | Firebase Cloud Messaging (FCM) |
| Real-time | WebSocket |
| Аутентификация | JWT (access + refresh tokens) |
| Документация | Swagger/OpenAPI |

---

## Архитектура

Два независимых приложения:

- **Admin API** (`cmd/admin`) — панель управления для системных пользователей
- **Mobile API** (`cmd/mobile`) — мобильное приложение для полевых сотрудников

---

## Роли и доступы

### Системные пользователи (Admin API)

| Роль | Описание | Доступ |
|---|---|---|
| `admin` | Администратор | Полный доступ ко всем функциям |
| `manager` | Менеджер | Мониторинг супервайзеров и мерчандайзеров, просмотр отчётов |
| `supervisor` | Супервайзер | Управление своей группой мерчандайзеров, создание/просмотр расписаний и визитов своей группы |

### Полевые сотрудники (Mobile API)

| Роль | Описание |
|---|---|
| `merchandiser` | Полевой сотрудник. Может быть в группе под управлением супервайзера |

> Супервайзеры могут быть как системными пользователями (в Admin API), так и привязаны к группе сотрудников через `supervisor_id` в модели Group.

---

## Группы и иерархия

Группы — центральная сущность для организации структуры как сотрудников, так и магазинов.

### Модель Group

```json
{
  "id": "uuid",
  "name": "Территория 1",
  "type": "staff | store",
  "parent_id": "uuid | null",
  "parent_name": "string | null",
  "supervisor_id": "uuid | null",
  "supervisor_name": "string | null",
  "color": "#FF5733",
  "layer": 1,
  "cluster": "cluster-name",
  "created_at": "timestamp",
  "groups": []
}
```

**Поля:**
- `type` — тип группы: `staff` (для сотрудников) или `store` (для магазинов)
- `parent_id` — ссылка на родительскую группу (поддержка иерархии группа → подгруппы)
- `supervisor_id` — прикреплённый супервайзер (для групп сотрудников)
- `color` — hex-цвет для отображения на карте (визуальное разделение подгрупп)
- `layer` — порядок слоя при отрисовке на карте
- `cluster` — имя кластера для группировки точек на карте
- `groups` — вложенные подгруппы (при запросе иерархии)

---

## Admin API

**Base URL:** `/api/v1/admin`

---

### Аутентификация системных пользователей

| Метод | Путь | Описание |
|---|---|---|
| `POST` | `/auth/system/login` | Вход в систему |
| `POST` | `/auth/system/register` | Регистрация системного пользователя |
| `POST` | `/auth/system/refresh` | Обновление access-токена |
| `POST` | `/auth/system/logout` | Выход |

---

### Системные пользователи

| Метод | Путь | Описание |
|---|---|---|
| `GET` | `/system-users` | Список системных пользователей |
| `POST` | `/system-users` | Создать системного пользователя |
| `PUT` | `/system-users/:id` | Обновить системного пользователя |
| `DELETE` | `/system-users/:id` | Удалить системного пользователя |

**CreateSystemUser (тело запроса):**
```json
{
  "full_name": "string (required)",
  "login": "string (required, alphanum)",
  "password": "string (required, 8-20 символов)",
  "confirm_password": "string (required)",
  "role": "admin | manager | supervisor (required)"
}
```

---

### Группы

#### Что нового: поддержка групп для магазинов и иерархии

Группы теперь поддерживают два типа (`staff` и `store`), а также вложенность через `parent_id`. Для карты мониторинга каждая группа имеет свой цвет (`color`), слой (`layer`) и кластер (`cluster`).

| Метод | Путь | Описание |
|---|---|---|
| `POST` | `/groups` | Создать группу (staff или store) |
| `GET` | `/groups` | Список групп с фильтрацией по типу |
| `GET` | `/groups/:id` | Получить группу по ID |
| `PUT` | `/groups/:id` | Обновить группу |
| `DELETE` | `/groups/:id` | Удалить группу |
| `GET` | `/groups/hierarchy/:type` | Получить иерархическую структуру групп |

**Query-параметры GET `/groups`:**
- `type` — `staff` или `store`
- `supervisor_id` — фильтр по супервайзеру

**GET `/groups/hierarchy/:type`** — возвращает дерево групп с вложенными подгруппами в поле `groups`.

**CreateGroup (тело запроса):**
```json
{
  "name": "Территория 1",
  "type": "staff | store",
  "parent_id": "uuid | null",
  "supervisor_id": "uuid | null",
  "color": "#FF5733",
  "layer": 1,
  "cluster": "north"
}
```

---

### Сотрудники (Staff)

| Метод | Путь | Описание |
|---|---|---|
| `GET` | `/staff` | Список сотрудников (с пагинацией, фильтрами) |
| `POST` | `/staff` | Создать сотрудника |
| `GET` | `/staff/:id` | Получить сотрудника по ID |
| `PUT` | `/staff/:id` | Обновить сотрудника |
| `PATCH` | `/staff/:id/deactivate` | Деактивировать сотрудника |
| `GET` | `/staff/online-offline` | Онлайн/офлайн статус всех сотрудников за период |
| `GET` | `/staff/workday/:id` | Рабочее время сотрудника |
| `GET` | `/staff/day-start/:id` | Проверить, начал ли сотрудник день |
| `PATCH` | `/staff/end-work/:staff_id` | Завершить рабочий день (admin) |

**Query-параметры GET `/staff`:**
- `group_id` — фильтр по группе
- `status` — фильтр по статусу
- `limit`, `offset` — пагинация

---

### Мониторинг

#### Что нового: фильтрация по группам на карте

| Метод | Путь | Описание |
|---|---|---|
| `GET` | `/staff/monitoring/real-time/:group_id` | Real-time локации сотрудников по группе |
| `GET` | `/staff/monitoring/track/:staff_id` | История перемещений сотрудника за дату |
| `GET` | `/admin/ws-location` | WebSocket — real-time обновления локаций |
| `GET` | `/admin/ws-notification` | WebSocket — уведомления в реальном времени |

При выборе группы в мониторинге (`/monitoring/real-time/:group_id`) возвращаются все сотрудники из группы. Визуальное разделение по подгруппам реализуется через поля `color`, `layer`, `cluster` модели Group.

---

### Магазины (Stores)

#### Что нового: привязка магазинов к группам типа `store`

| Метод | Путь | Описание |
|---|---|---|
| `GET` | `/stores` | Список магазинов с фильтрами |
| `GET` | `/stores/:id` | Получить магазин с фотографиями |
| `GET` | `/stores/by-group/:group_id` | Магазины по группе |
| `PATCH` | `/stores/:id/approve` | Подтвердить магазин |
| `PATCH` | `/stores/:id/group` | Привязать магазин к группе |
| `DELETE` | `/stores/:id/group` | Открепить магазин от группы |
| `DELETE` | `/stores/:id` | Удалить магазин |

**PATCH `/stores/:id/group` (тело запроса):**
```json
{
  "group_id": "uuid"
}
```

**Query-параметры GET `/stores`:**
- `is_active` — `true` / `false`
- `status` — статус магазина
- `sort_by` — поле сортировки
- `sort_order` — `asc` / `desc`

---

### Привязка магазинов к сотрудникам

| Метод | Путь | Описание |
|---|---|---|
| `GET` | `/staff-stores/:staff_id` | Магазины, привязанные к сотруднику |
| `POST` | `/staff-stores/:staff_id` | Привязать магазин к сотруднику |
| `DELETE` | `/staff-stores/:staff_id` | Открепить магазин от сотрудника |

---

### Расписания (Schedules)

#### Что нового: массовое создание расписаний через группу магазинов

| Метод | Путь | Описание |
|---|---|---|
| `GET` | `/staff/schedules` | Список всех расписаний (admin) |
| `POST` | `/staff/schedules/assign` | Назначить расписание одному магазину |
| `POST` | `/staff/schedule/assign-many` | Массовое создание расписаний |
| `PUT` | `/staff/schedules/:id` | Обновить расписание |
| `DELETE` | `/staff/schedules/:id` | Удалить расписание |

**Query-параметры GET `/staff/schedules`:**
- `staff_id` — фильтр по сотруднику
- `store_id` — фильтр по магазину
- `date_from`, `date_to` — диапазон дат
- `limit`, `offset` — пагинация

**POST `/staff/schedules/assign` (тело запроса):**
```json
{
  "staff_id": "uuid",
  "store_id": "uuid (required)",
  "planned_date": "2025-01-15T00:00:00Z (required)",
  "description": "string (required)",
  "window_open": "2025-01-15T09:00:00Z (required)",
  "window_close": "2025-01-15T18:00:00Z (required)",
  "is_photo_required": true,
  "is_selfie_required": false,
  "is_comment_required": true
}
```

**POST `/staff/schedule/assign-many` (тело запроса):**

Массовое создание расписаний — по одному на каждый магазин. Можно передать список магазинов напрямую (`store_ids`) или выбрать группу магазинов (`group_id`), тогда все магазины группы подтягиваются автоматически. Поддерживается частичное исключение (`exclude_store_ids`).

```json
{
  "staff_id": "uuid (required)",
  "group_id": "uuid | null",
  "store_ids": ["uuid", "uuid"],
  "exclude_store_ids": ["uuid"],
  "planned_date": "2025-01-15T00:00:00Z",
  "description": "string",
  "window_open": "2025-01-15T09:00:00Z",
  "window_close": "2025-01-15T18:00:00Z",
  "is_photo_required": true,
  "is_selfie_required": false,
  "is_comment_required": true
}
```

Ответ:
```json
{ "status": "success", "count": 12 }
```

---

### Визиты (Visits)

#### Что нового: массовое создание, учёт времени, статус проверки

| Метод | Путь | Описание |
|---|---|---|
| `POST` | `/visits/bulk` | Массовое создание визитов |
| `GET` | `/visits/review-queue` | Очередь визитов на проверку |
| `PATCH` | `/visits/:id/review` | Проверить визит (approve/reject/revision) |
| `GET` | `/reports/export` | Экспорт визитов в Excel |

**POST `/visits/bulk` (тело запроса):**

Массовое создание визитов — по одному на каждый магазин из расписания. Можно передать список магазинов или указать группу.

```json
{
  "schedule_id": "uuid (required)",
  "staff_id": "uuid (required)",
  "group_id": "uuid | null",
  "store_ids": ["uuid", "uuid"],
  "exclude_store_ids": ["uuid"],
  "lat": 41.2995,
  "lon": 69.2401
}
```

Ответ:
```json
{ "status": "success", "count": 8 }
```

**PATCH `/visits/:id/review` (тело запроса):**
```json
{
  "status": "approved | rejected | revision",
  "review_comment": "string (опционально)"
}
```

**Query-параметры GET `/visits/review-queue`:**
- `status` — фильтр по статусу проверки
- `limit`, `offset` — пагинация

**Query-параметры GET `/reports/export`:**
- `staff_id` — фильтр по сотруднику
- `store_id` — фильтр по магазину
- `begin`, `end` — диапазон дат (`YYYY-MM-DD`)

---

### Еженедельные расписания сотрудников

| Метод | Путь | Описание |
|---|---|---|
| `GET` | `/staff/time-schedule/:id` | Список недельных расписаний сотрудника |
| `PUT` | `/staff/time-schedule/:id` | Обновить недельное расписание |
| `DELETE` | `/staff/time-schedule/:id` | Удалить недельное расписание |
| `POST` | `/staff/schedules/assign` | Назначить недельное расписание |
| `POST` | `/staff/schedule/assign-many` | Массовое назначение |

---

### Уведомления

#### Что нового: уведомления о внеплановых визитах сотрудника

| Метод | Путь | Описание |
|---|---|---|
| `GET` | `/notifications` | Список уведомлений панели |
| `POST` | `/notifications` | Создать и отправить уведомление |
| `PATCH` | `/notifications/:id/read` | Отметить как прочитанное |

**POST `/notifications` (тело запроса):**
```json
{
  "user_id": "uuid | null",
  "group_id": "uuid | null",
  "title": "string",
  "body": "string",
  "type": "string",
  "meta": {}
}
```

> При создании внепланового визита сотрудником (через mobile API `/schedule`) — система автоматически создаёт push-уведомление для менеджера/супервайзера. В уведомлении передаётся `meta` с данными о сотруднике, магазине и геолокации.

---

### Аудит и настройки

| Метод | Путь | Описание |
|---|---|---|
| `GET` | `/audit/logs` | Журнал действий |
| `GET` | `/settings` | Список настроек системы |
| `POST` | `/settings` | Создать настройку |
| `PUT` | `/settings/:id` | Обновить настройку |
| `GET` | `/settings/:id` | Получить настройку |
| `DELETE` | `/settings/:id` | Удалить настройку |

**Query-параметры GET `/audit/logs`:**
- `action` — тип действия
- `target_table` — целевая таблица
- `date_from`, `date_to` — диапазон дат

---

### Хранилище

| Метод | Путь | Описание |
|---|---|---|
| `GET` | `/storage/upload-url` | Получить pre-signed URL для загрузки |
| `GET` | `/storage/download-url` | Получить pre-signed URL для скачивания |

---

## Mobile API

**Base URL:** `/api/v1`

---

### Аутентификация сотрудников

| Метод | Путь | Описание |
|---|---|---|
| `POST` | `/auth/staff/login` | Вход |
| `POST` | `/auth/staff/register` | Регистрация |
| `POST` | `/auth/staff/refresh` | Обновление токена |
| `POST` | `/auth/staff/logout` | Выход |

---

### Профиль

| Метод | Путь | Описание |
|---|---|---|
| `GET` | `/profile` | Профиль с визитами и расписаниями |
| `PUT` | `/profile` | Обновить профиль |
| `PUT` | `/profile/avatar` | Загрузить аватар |
| `PUT` | `/device-token` | Привязать/обновить FCM-токен устройства |

---

### Рабочий день

| Метод | Путь | Описание |
|---|---|---|
| `GET` | `/day-status` | Проверить, начат ли рабочий день |
| `PUT` | `/start-work` | Начать рабочий день |
| `PUT` | `/end-work` | Завершить рабочий день |

---

### Расписания

| Метод | Путь | Описание |
|---|---|---|
| `GET` | `/schedule` | Список назначенных расписаний |
| `POST` | `/schedule` | Создать внеплановый визит |

**POST `/schedule` — создание внепланового визита (тело запроса):**

Сотрудник может создать внеплановый визит самостоятельно (форс-мажор, изменение маршрута). Такой визит помечается флагом `is_unplanned = true` и фиксирует геолокацию создания.

```json
{
  "store_id": "uuid (required)",
  "planned_date": "2025-01-15T00:00:00Z (required)",
  "description": "Причина или комментарий (required)",
  "window_open": "2025-01-15T09:00:00Z (required)",
  "window_close": "2025-01-15T18:00:00Z (required)",
  "is_photo_required": false,
  "is_selfie_required": false,
  "is_comment_required": false,
  "lat": 41.2995,
  "lon": 69.2401
}
```

Ответ:
```json
{ "message": "unplanned schedule created" }
```

> После создания внепланового расписания система отправляет push-уведомление администратору/менеджеру для одобрения или отклонения.

---

### Визиты

#### Что нового: учёт времени нахождения на точке

| Метод | Путь | Описание |
|---|---|---|
| `POST` | `/visits/check-in` | Отметиться на точке (check-in) |
| `GET` | `/visits` | Список визитов по диапазону дат |
| `GET` | `/visits/:id` | Детали визита с фотографиями |
| `POST` | `/visits/:id/photo` | Загрузить фото визита |
| `POST` | `/visits/:id/complete` | Завершить визит (check-out) |
| `PUT` | `/visits/:id` | Обновить визит |
| `DELETE` | `/visits/:id/photo` | Удалить фото визита |

**POST `/visits/check-in` (тело запроса):**
```json
{
  "schedule_id": "uuid (required)",
  "store_id": "uuid (required)",
  "lat": 41.2995,
  "lon": 69.2401
}
```

Фиксирует `check_in_at` и координаты прибытия. Возвращает созданный визит.

**POST `/visits/:id/complete` (тело запроса):**

Завершение визита фиксирует `check_out_at`. Поле `duration` рассчитывается автоматически как разница между `check_out_at` и `check_in_at` в секундах.

```json
{
  "check_out_at": "2025-01-15T11:30:00Z",
  "comment": "string (если требуется)",
  "is_offline": false
}
```

**Модель Visit:**
```json
{
  "id": "uuid",
  "schedule_id": "uuid | null",
  "staff_id": "uuid",
  "store_id": "uuid",
  "check_in_at": "2025-01-15T09:15:00Z",
  "check_out_at": "2025-01-15T11:30:00Z",
  "check_in_lat": 41.2995,
  "check_in_lon": 69.2401,
  "distance_deviation": 150,
  "duration": 8100,
  "comment": "string | null",
  "is_offline": false,
  "review_status": "pending | approved | rejected | revision",
  "review_comment": "string | null",
  "reviewer_id": "uuid | null",
  "reviewed_at": "timestamp | null",
  "created_at": "timestamp"
}
```

**Поля учёта времени:**
- `check_in_at` — время прибытия на точку
- `check_out_at` — время ухода с точки
- `duration` — длительность визита в секундах (рассчитывается автоматически)
- `distance_deviation` — отклонение от плановой точки в метрах

**Query-параметры GET `/visits`:**
- `begin`, `end` — диапазон дат (`YYYY-MM-DD`, оба обязательны)
- `limit`, `offset` — пагинация

---

### Магазины

| Метод | Путь | Описание |
|---|---|---|
| `GET` | `/my-stores` | Магазины, назначенные сотруднику |
| `POST` | `/stores/upload-photo` | Загрузить фото магазина |
| `POST` | `/stores/propose` | Предложить новый магазин |

---

### Геолокация

| Метод | Путь | Описание |
|---|---|---|
| `POST` | `/location` | Отправить пакет координат |

**POST `/location` (тело запроса):**
```json
[
  {
    "lat": 41.2995,
    "lon": 69.2401,
    "battery_level": 85,
    "captured_at": "2025-01-15T10:00:00Z",
    "status": "moving | idle"
  }
]
```

---

### Уведомления

| Метод | Путь | Описание |
|---|---|---|
| `GET` | `/notifications` | Список уведомлений сотрудника |
| `PATCH` | `/notifications/:id/read` | Отметить как прочитанное |

---

### Прочее

| Метод | Путь | Описание |
|---|---|---|
| `GET` | `/settings` | Настройки системы |
| `GET` | `/storage/upload-url` | Pre-signed URL для загрузки |
| `GET` | `/storage/download-url` | Pre-signed URL для скачивания |
| `GET` | `/ping` | Health check |

---

## Сводка изменений по задачам

### 1. Мониторинг → Группы и подгруппы
- Эндпоинт `GET /staff/monitoring/real-time/:group_id` принимает ID группы
- `GET /groups/hierarchy/:type` возвращает полное дерево групп с вложенными подгруппами
- Модель `Group` содержит `color`, `layer`, `cluster` для визуального разделения на карте

### 2. Группы для магазинов
- `POST /groups` с `"type": "store"` — создание группы магазинов
- `GET /groups?type=store` — список групп магазинов
- `PATCH /stores/:id/group` — привязать магазин к группе
- `DELETE /stores/:id/group` — открепить магазин от группы
- `GET /stores/by-group/:group_id` — магазины по группе

### 3. Массовое создание визитов
- `POST /visits/bulk` — массовое создание визитов
- Поддержка `group_id` (подтягивает все магазины группы) и `store_ids` (явный список)
- `exclude_store_ids` — частичное исключение магазинов из списка

### 4. Создание визита сотрудником (внеплановый)
- `POST /schedule` (Mobile API) — сотрудник создаёт внеплановый визит
- Поле `is_unplanned = true` автоматически проставляется
- Фиксируется геолокация (`created_at_lat`, `created_at_lon`) и время создания
- Система отправляет push-уведомление менеджеру/супервайзеру

### 5. Учёт времени нахождения на точке
- `POST /visits/check-in` — фиксирует `check_in_at` + координаты
- `POST /visits/:id/complete` — фиксирует `check_out_at`
- Поле `duration` (секунды) рассчитывается автоматически
- Доступно в карточке визита и в экспорте Excel (`GET /reports/export`)

### 6. Пользователи и доступы
- Системные пользователи: `admin`, `manager`, `supervisor`
- Полевые сотрудники: `merchandiser`
- Создание через `POST /system-users` с указанием роли
- Супервайзер привязывается к группе сотрудников через `supervisor_id` в модели Group

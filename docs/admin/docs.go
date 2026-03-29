package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "swagger": "2.0",
    "info": {
        "title": "GPS Tracker — Admin API",
        "description": "Административная панель GPS-трекера: управление персоналом, магазинами, группами, расписаниями, визитами и настройками.",
        "version": "1.0",
        "contact": {}
    },
    "basePath": "/api/v1",
    "schemes": ["http", "https"],
    "securityDefinitions": {
        "BearerAuth": {
            "type": "apiKey",
            "name": "Authorization",
            "in": "header",
            "description": "Введите токен в формате: Bearer {token}"
        }
    },
    "tags": [
        {"name": "Health", "description": "Проверка состояния сервиса"},
        {"name": "Storage", "description": "Работа с файлами (MinIO)"},
        {"name": "Auth", "description": "Аутентификация системных пользователей"},
        {"name": "Staff Users", "description": "Управление полевыми сотрудниками"},
        {"name": "Monitoring", "description": "Мониторинг местоположения сотрудников"},
        {"name": "Stores", "description": "Управление торговыми точками"},
        {"name": "Groups", "description": "Управление группами"},
        {"name": "Staff Mapping", "description": "Привязка сотрудников к магазинам"},
        {"name": "Visits", "description": "Визиты и отчёты"},
        {"name": "Settings", "description": "Системные настройки"},
        {"name": "Audit", "description": "Журнал аудита"},
        {"name": "System Users", "description": "Управление системными пользователями (admin/supervisor)"},
        {"name": "Notifications", "description": "Системные уведомления для admin/supervisor"}
    ],
    "paths": {
        "/ping": {
            "get": {
                "tags": ["Health"],
                "summary": "Health-check",
                "produces": ["application/json"],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "object",
                            "properties": {
                                "message": {"type": "string", "example": "Pong!"}
                            }
                        }
                    }
                }
            }
        },
        "/storage/download-url": {
            "get": {
                "tags": ["Storage"],
                "summary": "Получить presigned URL для скачивания файла",
                "produces": ["application/json"],
                "parameters": [
                    {"name": "file_name", "in": "query", "required": true, "type": "string", "description": "Имя файла в хранилище"}
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "object",
                            "properties": {
                                "download_url": {"type": "string"}
                            }
                        }
                    },
                    "400": {"description": "file_name is required"},
                    "500": {"description": "Internal Server Error"}
                }
            }
        },
        "/storage/upload-url": {
            "get": {
                "tags": ["Storage"],
                "summary": "Получить presigned URL для загрузки файла",
                "produces": ["application/json"],
                "parameters": [
                    {"name": "file_name", "in": "query", "required": true, "type": "string", "description": "Желаемое имя файла"}
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "object",
                            "properties": {
                                "upload_url": {"type": "string"}
                            }
                        }
                    },
                    "500": {"description": "Internal Server Error"}
                }
            }
        },
        "/auth/system/login": {
            "post": {
                "tags": ["Auth"],
                "summary": "Вход в систему (admin/supervisor)",
                "consumes": ["application/json"],
                "produces": ["application/json"],
                "parameters": [
                    {
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {"$ref": "#/definitions/Auth"}
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Успешная авторизация",
                        "schema": {"$ref": "#/definitions/TokenPair"}
                    },
                    "401": {"description": "Неверные учётные данные"},
                    "500": {"description": "Internal Server Error"}
                }
            }
        },
        "/auth/system/register": {
            "post": {
                "tags": ["Auth"],
                "summary": "Регистрация системного пользователя",
                "consumes": ["application/json"],
                "produces": ["application/json"],
                "parameters": [
                    {
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {"$ref": "#/definitions/CreateUser"}
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Пользователь создан",
                        "schema": {
                            "type": "object",
                            "properties": {
                                "ok": {"type": "string", "example": "created user"}
                            }
                        }
                    },
                    "400": {"description": "Пользователь уже существует или ошибка валидации"},
                    "500": {"description": "Internal Server Error"}
                }
            }
        },
        "/auth/system/refresh": {
            "post": {
                "tags": ["Auth"],
                "summary": "Обновить пару токенов",
                "consumes": ["application/json"],
                "produces": ["application/json"],
                "parameters": [
                    {
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {"$ref": "#/definitions/RefreshRequest"}
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Новая пара токенов",
                        "schema": {"$ref": "#/definitions/TokenPair"}
                    },
                    "401": {"description": "Невалидный refresh token"}
                }
            }
        },
        "/auth/system/logout": {
            "post": {
                "tags": ["Auth"],
                "summary": "Выход из системы",
                "consumes": ["application/json"],
                "produces": ["application/json"],
                "parameters": [
                    {
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {"$ref": "#/definitions/RefreshRequest"}
                    }
                ],
                "responses": {
                    "200": {"description": "OK"},
                    "400": {"description": "Refresh token не найден"}
                }
            }
        },
        "/admin/reports/export": {
            "get": {
                "tags": ["Visits"],
                "summary": "Экспорт отчёта по визитам в Excel",
                "security": [{"BearerAuth": []}],
                "produces": ["application/octet-stream"],
                "parameters": [
                    {"name": "staff_id", "in": "query", "type": "string", "format": "uuid", "description": "Фильтр по сотруднику (опционально)"},
                    {"name": "store_id", "in": "query", "type": "string", "format": "uuid", "description": "Фильтр по магазину (опционально)"},
                    {"name": "begin", "in": "query", "type": "string", "format": "date", "description": "Дата начала (YYYY-MM-DD, опционально)"},
                    {"name": "end", "in": "query", "type": "string", "format": "date", "description": "Дата окончания (YYYY-MM-DD, опционально)"}
                ],
                "responses": {
                    "200": {"description": "Excel-файл с отчётом"},
                    "400": {"description": "Невалидные параметры"},
                    "500": {"description": "Internal Server Error"}
                }
            }
        },
        "/admin/audit/logs": {
            "get": {
                "tags": ["Audit"],
                "summary": "Получить журнал аудита",
                "security": [{"BearerAuth": []}],
                "produces": ["application/json"],
                "responses": {
                    "200": {
                        "description": "Список записей аудита",
                        "schema": {
                            "type": "object",
                            "properties": {
                                "audit_logs": {
                                    "type": "array",
                                    "items": {"$ref": "#/definitions/AuditLog"}
                                }
                            }
                        }
                    },
                    "500": {"description": "Internal Server Error"}
                }
            }
        },
        "/admin/staff": {
            "get": {
                "tags": ["Staff Users"],
                "summary": "Список всех сотрудников",
                "security": [{"BearerAuth": []}],
                "produces": ["application/json"],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "object",
                            "properties": {
                                "staff_users": {
                                    "type": "array",
                                    "items": {"$ref": "#/definitions/StaffUser"}
                                }
                            }
                        }
                    },
                    "500": {"description": "Internal Server Error"}
                }
            },
            "post": {
                "tags": ["Staff Users"],
                "summary": "Создать сотрудника",
                "security": [{"BearerAuth": []}],
                "consumes": ["application/json"],
                "produces": ["application/json"],
                "parameters": [
                    {
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {"$ref": "#/definitions/CreateStaffUser"}
                    }
                ],
                "responses": {
                    "201": {"description": "Сотрудник создан"},
                    "400": {"description": "Ошибка валидации"},
                    "403": {"description": "Нет доступа"},
                    "500": {"description": "Internal Server Error"}
                }
            }
        },
        "/admin/staff/{id}": {
            "get": {
                "tags": ["Staff Users"],
                "summary": "Получить сотрудника по ID",
                "security": [{"BearerAuth": []}],
                "produces": ["application/json"],
                "parameters": [
                    {"name": "id", "in": "path", "required": true, "type": "string", "format": "uuid"}
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "object",
                            "properties": {
                                "staff_user": {"$ref": "#/definitions/StaffUser"}
                            }
                        }
                    },
                    "400": {"description": "Невалидный ID или пользователь не найден"},
                    "500": {"description": "Internal Server Error"}
                }
            },
            "put": {
                "tags": ["Staff Users"],
                "summary": "Обновить данные сотрудника",
                "security": [{"BearerAuth": []}],
                "consumes": ["application/json"],
                "produces": ["application/json"],
                "parameters": [
                    {"name": "id", "in": "path", "required": true, "type": "string", "format": "uuid"},
                    {
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {"$ref": "#/definitions/UpdateStaffUser"}
                    }
                ],
                "responses": {
                    "200": {"description": "Обновлено"},
                    "400": {"description": "Ошибка валидации"},
                    "403": {"description": "Нет доступа"},
                    "404": {"description": "Сотрудник не найден"},
                    "500": {"description": "Internal Server Error"}
                }
            }
        },
        "/admin/staff/{id}/deactivate": {
            "patch": {
                "tags": ["Staff Users"],
                "summary": "Деактивировать сотрудника",
                "security": [{"BearerAuth": []}],
                "consumes": ["application/json"],
                "produces": ["application/json"],
                "parameters": [
                    {"name": "id", "in": "path", "required": true, "type": "string", "format": "uuid"},
                    {
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {"$ref": "#/definitions/DeactivateStaffUser"}
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Деактивирован",
                        "schema": {
                            "type": "object",
                            "properties": {
                                "message": {"type": "string", "example": "deactivation complete"}
                            }
                        }
                    },
                    "400": {"description": "Ошибка валидации"},
                    "500": {"description": "Internal Server Error"}
                }
            }
        },
        "/admin/staff/schedules": {
            "get": {
                "tags": ["Staff Users"],
                "summary": "Список расписаний (с фильтрацией)",
                "security": [{"BearerAuth": []}],
                "produces": ["application/json"],
                "parameters": [
                    {"name": "limit", "in": "query", "type": "integer", "default": 20},
                    {"name": "offset", "in": "query", "type": "integer", "default": 0},
                    {"name": "staff_id", "in": "query", "type": "string", "format": "uuid", "description": "Фильтр по сотруднику"},
                    {"name": "store_id", "in": "query", "type": "string", "format": "uuid", "description": "Фильтр по магазину"},
                    {"name": "date_from", "in": "query", "type": "string", "format": "date", "description": "Дата начала (YYYY-MM-DD)"},
                    {"name": "date_to", "in": "query", "type": "string", "format": "date", "description": "Дата окончания (YYYY-MM-DD)"}
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "object",
                            "properties": {
                                "data": {
                                    "type": "array",
                                    "items": {"$ref": "#/definitions/Schedule"}
                                },
                                "pagination": {"$ref": "#/definitions/Pagination"}
                            }
                        }
                    },
                    "500": {"description": "Internal Server Error"}
                }
            }
        },
        "/admin/staff/schedules/assign": {
            "post": {
                "tags": ["Staff Users"],
                "summary": "Назначить расписание сотруднику",
                "security": [{"BearerAuth": []}],
                "consumes": ["application/json"],
                "produces": ["application/json"],
                "parameters": [
                    {
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {"$ref": "#/definitions/CreateSchedule"}
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Расписание создано",
                        "schema": {
                            "type": "object",
                            "properties": {
                                "schedule": {"$ref": "#/definitions/CreateSchedule"}
                            }
                        }
                    },
                    "403": {"description": "Нет доступа"},
                    "500": {"description": "Internal Server Error"}
                }
            }
        },
        "/admin/staff/schedules/{id}": {
            "delete": {
                "tags": ["Staff Users"],
                "summary": "Удалить назначение расписания (schedule) сотруднику",
                "security": [{"BearerAuth": []}],
                "produces": ["application/json"],
                "parameters": [
                    {"name": "id", "in": "path", "required": true, "type": "string", "format": "uuid", "description": "ID расписания"}
                ],
                "responses": {
                    "200": {"description": "Удалено"},
                    "400": {"description": "Невалидный ID"},
                    "500": {"description": "Internal Server Error"}
                }
            }
        },
        "/admin/staff/time-schedule/{id}": {
            "get": {
                "tags": ["Staff Users"],
                "summary": "Получить time-schedule записи сотрудника",
                "description": "Возвращает список временны́х расписаний (рабочие дни/часы) для указанного сотрудника.",
                "security": [{"BearerAuth": []}],
                "produces": ["application/json"],
                "parameters": [
                    {"name": "id", "in": "path", "required": true, "type": "string", "format": "uuid", "description": "ID сотрудника"}
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "object",
                            "properties": {
                                "schedules": {
                                    "type": "array",
                                    "items": {"$ref": "#/definitions/StaffUserSchedule"}
                                }
                            }
                        }
                    },
                    "400": {"description": "Невалидный ID"},
                    "500": {"description": "Internal Server Error"}
                }
            },
            "put": {
                "tags": ["Staff Users"],
                "summary": "Обновить time-schedule запись сотрудника",
                "description": "Частичное обновление временного расписания (рабочий день, начало/конец смены).",
                "security": [{"BearerAuth": []}],
                "consumes": ["application/json"],
                "produces": ["application/json"],
                "parameters": [
                    {"name": "id", "in": "path", "required": true, "type": "string", "format": "uuid", "description": "ID time-schedule записи"},
                    {
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {"$ref": "#/definitions/UpdateStaffSchedule"}
                    }
                ],
                "responses": {
                    "200": {"description": "Обновлено"},
                    "400": {"description": "Ошибка валидации"},
                    "500": {"description": "Internal Server Error"}
                }
            },
            "delete": {
                "tags": ["Staff Users"],
                "summary": "Удалить time-schedule запись сотрудника",
                "security": [{"BearerAuth": []}],
                "produces": ["application/json"],
                "parameters": [
                    {"name": "id", "in": "path", "required": true, "type": "string", "format": "uuid", "description": "ID time-schedule записи"}
                ],
                "responses": {
                    "200": {"description": "Удалено"},
                    "400": {"description": "Невалидный ID"},
                    "500": {"description": "Internal Server Error"}
                }
            }
        },
        "/admin/staff/online-offline": {
            "get": {
                "tags": ["Staff Users"],
                "summary": "Онлайн/оффлайн статус всех сотрудников",
                "security": [{"BearerAuth": []}],
                "produces": ["application/json"],
                "parameters": [
                    {"name": "begin", "in": "query", "required": true, "type": "string", "format": "date", "description": "YYYY-MM-DD"},
                    {"name": "end", "in": "query", "required": true, "type": "string", "format": "date", "description": "YYYY-MM-DD"}
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "object",
                            "properties": {
                                "data": {
                                    "type": "array",
                                    "items": {
                                        "type": "object",
                                        "properties": {
                                            "staff_id": {"type": "string", "format": "uuid"},
                                            "online": {"type": "number", "description": "Минуты онлайн"},
                                            "offline": {"type": "number", "description": "Минуты оффлайн"}
                                        }
                                    }
                                }
                            }
                        }
                    },
                    "500": {"description": "Internal Server Error"}
                }
            }
        },
        "/admin/staff/workday/{id}": {
            "get": {
                "tags": ["Staff Users"],
                "summary": "Получить рабочее время сотрудника",
                "security": [{"BearerAuth": []}],
                "produces": ["application/json"],
                "parameters": [
                    {"name": "id", "in": "path", "required": true, "type": "string", "format": "uuid"}
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "object",
                            "properties": {
                                "time": {"type": "object"}
                            }
                        }
                    },
                    "400": {"description": "Невалидный ID"},
                    "500": {"description": "Internal Server Error"}
                }
            }
        },
        "/admin/staff/day-start/{id}": {
            "get": {
                "tags": ["Staff Users"],
                "summary": "Проверить, начал ли сотрудник рабочий день",
                "security": [{"BearerAuth": []}],
                "produces": ["application/json"],
                "parameters": [
                    {"name": "id", "in": "path", "required": true, "type": "string", "format": "uuid"}
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "object",
                            "properties": {
                                "day_started": {"type": "boolean"}
                            }
                        }
                    },
                    "400": {"description": "Невалидный ID"},
                    "500": {"description": "Internal Server Error"}
                }
            }
        },
        "/admin/staff/monitoring/real-time": {
            "get": {
                "tags": ["Monitoring"],
                "summary": "Текущее местоположение сотрудников группы",
                "security": [{"BearerAuth": []}],
                "produces": ["application/json"],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "object",
                            "properties": {
                                "staff_users": {
                                    "type": "array",
                                    "items": {"$ref": "#/definitions/StaffUserLocation"}
                                }
                            }
                        }
                    },
                    "400": {"description": "Ошибка"},
                    "500": {"description": "Internal Server Error"}
                }
            }
        },
        "/admin/staff/monitoring/track/{staff_id}": {
            "get": {
                "tags": ["Monitoring"],
                "summary": "История перемещений сотрудника за период",
                "security": [{"BearerAuth": []}],
                "produces": ["application/json"],
                "parameters": [
                    {"name": "staff_id", "in": "path", "required": true, "type": "string", "format": "uuid"},
                    {"name": "begin", "in": "query", "required": true, "type": "string", "format": "date", "description": "Дата начала (YYYY-MM-DD)"},
                    {"name": "end", "in": "query", "required": true, "type": "string", "format": "date", "description": "Дата окончания (YYYY-MM-DD)"}
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "object",
                            "properties": {
                                "staff_user": {"$ref": "#/definitions/StaffUserLocationHistory"}
                            }
                        }
                    },
                    "400": {"description": "Невалидные параметры"}
                }
            }
        },
        "/admin/stores": {
            "get": {
                "tags": ["Stores"],
                "summary": "Список магазинов с фильтрацией",
                "security": [{"BearerAuth": []}],
                "consumes": ["application/json"],
                "produces": ["application/json"],
                "parameters": [
                    {
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {"$ref": "#/definitions/StoreListFilter"}
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "object",
                            "properties": {
                                "stores": {
                                    "type": "array",
                                    "items": {"$ref": "#/definitions/Store"}
                                },
                                "total": {"type": "integer"}
                            }
                        }
                    },
                    "500": {"description": "Internal Server Error"}
                }
            }
        },
        "/admin/stores/{id}": {
            "get": {
                "tags": ["Stores"],
                "summary": "Получить магазин по ID",
                "security": [{"BearerAuth": []}],
                "produces": ["application/json"],
                "parameters": [
                    {"name": "id", "in": "path", "required": true, "type": "string", "format": "uuid"}
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "object",
                            "properties": {
                                "store": {"$ref": "#/definitions/Store"}
                            }
                        }
                    },
                    "400": {"description": "Невалидный ID"},
                    "404": {"description": "Магазин не найден"},
                    "500": {"description": "Internal Server Error"}
                }
            },
            "delete": {
                "tags": ["Stores"],
                "summary": "Удалить магазин",
                "security": [{"BearerAuth": []}],
                "produces": ["application/json"],
                "parameters": [
                    {"name": "id", "in": "path", "required": true, "type": "string", "format": "uuid", "description": "ID магазина"}
                ],
                "responses": {
                    "200": {"description": "Магазин удалён"},
                    "400": {"description": "Невалидный ID"},
                    "500": {"description": "Internal Server Error"}
                }
            }
        },
        "/admin/stores/{id}/approve": {
            "patch": {
                "tags": ["Stores"],
                "summary": "Подтвердить / обновить магазин",
                "security": [{"BearerAuth": []}],
                "consumes": ["application/json"],
                "produces": ["application/json"],
                "parameters": [
                    {"name": "id", "in": "path", "required": true, "type": "string", "format": "uuid"},
                    {
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {"$ref": "#/definitions/UpdateStore"}
                    }
                ],
                "responses": {
                    "200": {"description": "Обновлено"},
                    "404": {"description": "Магазин не найден"},
                    "500": {"description": "Internal Server Error"}
                }
            }
        },
        "/admin/groups": {
            "get": {
                "tags": ["Groups"],
                "summary": "Список групп",
                "security": [{"BearerAuth": []}],
                "produces": ["application/json"],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "object",
                            "properties": {
                                "groups": {
                                    "type": "array",
                                    "items": {"$ref": "#/definitions/Group"}
                                }
                            }
                        }
                    },
                    "500": {"description": "Internal Server Error"}
                }
            },
            "post": {
                "tags": ["Groups"],
                "summary": "Создать группу",
                "security": [{"BearerAuth": []}],
                "consumes": ["application/json"],
                "produces": ["application/json"],
                "parameters": [
                    {
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {"$ref": "#/definitions/Group"}
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Группа создана",
                        "schema": {
                            "type": "object",
                            "properties": {
                                "message": {"type": "string"}
                            }
                        }
                    },
                    "500": {"description": "Internal Server Error"}
                }
            }
        },
        "/admin/groups/{id}": {
            "get": {
                "tags": ["Groups"],
                "summary": "Получить группу по ID",
                "security": [{"BearerAuth": []}],
                "produces": ["application/json"],
                "parameters": [
                    {"name": "id", "in": "path", "required": true, "type": "string", "format": "uuid"}
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "object",
                            "properties": {
                                "group": {"$ref": "#/definitions/Group"}
                            }
                        }
                    },
                    "400": {"description": "Группа не найдена"}
                }
            },
            "put": {
                "tags": ["Groups"],
                "summary": "Обновить группу",
                "security": [{"BearerAuth": []}],
                "consumes": ["application/json"],
                "produces": ["application/json"],
                "parameters": [
                    {"name": "id", "in": "path", "required": true, "type": "string", "format": "uuid"},
                    {
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {"$ref": "#/definitions/Group"}
                    }
                ],
                "responses": {
                    "200": {"description": "Обновлено"},
                    "400": {"description": "Группа не найдена"},
                    "500": {"description": "Internal Server Error"}
                }
            },
            "delete": {
                "tags": ["Groups"],
                "summary": "Удалить группу",
                "security": [{"BearerAuth": []}],
                "produces": ["application/json"],
                "parameters": [
                    {"name": "id", "in": "path", "required": true, "type": "string", "format": "uuid"}
                ],
                "responses": {
                    "200": {"description": "Удалено"},
                    "400": {"description": "Группа не найдена"},
                    "500": {"description": "Internal Server Error"}
                }
            }
        },
        "/admin/staff-mapping/{staff_id}": {
            "get": {
                "tags": ["Staff Mapping"],
                "summary": "Список магазинов, назначенных сотруднику через маппинг",
                "description": "Возвращает магазины, привязанные к сотруднику через таблицу staff_stores_mapping.",
                "security": [{"BearerAuth": []}],
                "produces": ["application/json"],
                "parameters": [
                    {"name": "staff_id", "in": "path", "required": true, "type": "string", "format": "uuid", "description": "ID сотрудника"},
                    {"name": "limit", "in": "query", "required": true, "type": "integer", "description": "Количество записей"},
                    {"name": "offset", "in": "query", "required": true, "type": "integer", "description": "Смещение"}
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "object",
                            "properties": {
                                "data": {
                                    "type": "array",
                                    "items": {"$ref": "#/definitions/Store"}
                                },
                                "total": {"type": "integer"},
                                "limit": {"type": "integer"},
                                "offset": {"type": "integer"}
                            }
                        }
                    },
                    "400": {"description": "Неверный staff_id"},
                    "403": {"description": "Нет доступа"},
                    "500": {"description": "Internal Server Error"}
                }
            },
            "post": {
                "tags": ["Staff Mapping"],
                "summary": "Привязать магазин к сотруднику",
                "description": "Создаёт запись в staff_stores_mapping — назначает магазин сотруднику.",
                "security": [{"BearerAuth": []}],
                "consumes": ["application/json"],
                "produces": ["application/json"],
                "parameters": [
                    {"name": "staff_id", "in": "path", "required": true, "type": "string", "format": "uuid", "description": "ID сотрудника"},
                    {
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "type": "object",
                            "required": ["store_id"],
                            "properties": {
                                "store_id": {"type": "string", "format": "uuid", "description": "ID магазина"}
                            }
                        }
                    }
                ],
                "responses": {
                    "201": {"description": "Привязка создана"},
                    "400": {"description": "Неверные параметры"},
                    "403": {"description": "Нет доступа"},
                    "500": {"description": "Internal Server Error"}
                }
            },
            "delete": {
                "tags": ["Staff Mapping"],
                "summary": "Отвязать магазин от сотрудника",
                "description": "Удаляет запись из staff_stores_mapping по store_id и staff_id.",
                "security": [{"BearerAuth": []}],
                "consumes": ["application/json"],
                "produces": ["application/json"],
                "parameters": [
                    {"name": "staff_id", "in": "path", "required": true, "type": "string", "format": "uuid", "description": "ID сотрудника"},
                    {
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "type": "object",
                            "required": ["store_id"],
                            "properties": {
                                "store_id": {"type": "string", "format": "uuid", "description": "ID магазина"}
                            }
                        }
                    }
                ],
                "responses": {
                    "201": {"description": "Привязка удалена"},
                    "400": {"description": "Неверные параметры"},
                    "403": {"description": "Нет доступа"},
                    "404": {"description": "Привязка не найдена"},
                    "500": {"description": "Internal Server Error"}
                }
            }
        },
        "/admin/visit/review-queue": {
            "get": {
                "tags": ["Visits"],
                "summary": "Очередь визитов на проверку",
                "security": [{"BearerAuth": []}],
                "produces": ["application/json"],
                "parameters": [
                    {"name": "limit", "in": "query", "required": true, "type": "integer"},
                    {"name": "offset", "in": "query", "required": true, "type": "integer"},
                    {"name": "status", "in": "query", "type": "string", "enum": ["pending", "approved", "rejected", "revision"], "description": "Фильтр по статусу проверки"}
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "object",
                            "properties": {
                                "visits": {
                                    "type": "array",
                                    "items": {"$ref": "#/definitions/Visit"}
                                },
                                "total": {"type": "integer"}
                            }
                        }
                    },
                    "403": {"description": "Нет доступа"},
                    "500": {"description": "Internal Server Error"}
                }
            }
        },
        "/admin/visit/{id}/review": {
            "patch": {
                "tags": ["Visits"],
                "summary": "Рассмотреть визит (approve / reject)",
                "security": [{"BearerAuth": []}],
                "consumes": ["application/json"],
                "produces": ["application/json"],
                "parameters": [
                    {"name": "id", "in": "path", "required": true, "type": "string", "format": "uuid"},
                    {
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {"$ref": "#/definitions/ReviewVisitRequest"}
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Рассмотрено",
                        "schema": {
                            "type": "object",
                            "properties": {
                                "message": {"type": "string"},
                                "status": {"type": "string"}
                            }
                        }
                    },
                    "400": {"description": "status must be 'approved', 'rejected' or 'revision'"},
                    "409": {"description": "Визит уже имеет финальный статус (approved/rejected) и не может быть изменён"},
                    "403": {"description": "Нет доступа"},
                    "404": {"description": "Визит не найден"},
                    "500": {"description": "Internal Server Error"}
                }
            }
        },
        "/admin/settings": {
            "get": {
                "tags": ["Settings"],
                "summary": "Список системных настроек",
                "security": [{"BearerAuth": []}],
                "produces": ["application/json"],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "object",
                            "properties": {
                                "system_settings": {
                                    "type": "array",
                                    "items": {"$ref": "#/definitions/SystemSettings"}
                                }
                            }
                        }
                    },
                    "500": {"description": "Internal Server Error"}
                }
            },
            "post": {
                "tags": ["Settings"],
                "summary": "Создать системную настройку",
                "security": [{"BearerAuth": []}],
                "consumes": ["application/json"],
                "produces": ["application/json"],
                "parameters": [
                    {
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {"$ref": "#/definitions/SystemSettings"}
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Создано",
                        "schema": {
                            "type": "object",
                            "properties": {
                                "message": {"type": "string"}
                            }
                        }
                    },
                    "500": {"description": "Internal Server Error"}
                }
            }
        },
        "/admin/settings/{id}": {
            "get": {
                "tags": ["Settings"],
                "summary": "Получить настройку по ID",
                "security": [{"BearerAuth": []}],
                "produces": ["application/json"],
                "parameters": [
                    {"name": "id", "in": "path", "required": true, "type": "string", "format": "uuid"}
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "object",
                            "properties": {
                                "setting": {"$ref": "#/definitions/SystemSettings"}
                            }
                        }
                    },
                    "400": {"description": "Не найдено"},
                    "500": {"description": "Internal Server Error"}
                }
            },
            "put": {
                "tags": ["Settings"],
                "summary": "Обновить системную настройку",
                "security": [{"BearerAuth": []}],
                "consumes": ["application/json"],
                "produces": ["application/json"],
                "parameters": [
                    {"name": "id", "in": "path", "required": true, "type": "string", "format": "uuid"},
                    {
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {"$ref": "#/definitions/UpdateSystemSettings"}
                    }
                ],
                "responses": {
                    "200": {"description": "Обновлено"},
                    "500": {"description": "Internal Server Error"}
                }
            },
            "delete": {
                "tags": ["Settings"],
                "summary": "Удалить системную настройку",
                "security": [{"BearerAuth": []}],
                "produces": ["application/json"],
                "parameters": [
                    {"name": "id", "in": "path", "required": true, "type": "string", "format": "uuid"}
                ],
                "responses": {
                    "200": {"description": "Удалено"},
                    "500": {"description": "Internal Server Error"}
                }
            }
        },
        "/admin/system": {
            "get": {
                "tags": ["System Users"],
                "summary": "Список системных пользователей",
                "security": [{"BearerAuth": []}],
                "produces": ["application/json"],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "object",
                            "properties": {
                                "system_users": {
                                    "type": "array",
                                    "items": {"$ref": "#/definitions/SystemUser"}
                                }
                            }
                        }
                    },
                    "403": {"description": "Нет доступа (только admin)"},
                    "500": {"description": "Internal Server Error"}
                }
            },
            "post": {
                "tags": ["System Users"],
                "summary": "Создать системного пользователя",
                "security": [{"BearerAuth": []}],
                "consumes": ["application/json"],
                "produces": ["application/json"],
                "parameters": [
                    {
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {"$ref": "#/definitions/CreateSystemUser"}
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Создано",
                        "schema": {
                            "type": "object",
                            "properties": {
                                "message": {"type": "string", "example": "created"}
                            }
                        }
                    },
                    "400": {"description": "Ошибка валидации"},
                    "403": {"description": "Нет доступа (только admin)"},
                    "500": {"description": "Internal Server Error"}
                }
            }
        },
        "/admin/system/{id}": {
            "put": {
                "tags": ["System Users"],
                "summary": "Обновить системного пользователя",
                "security": [{"BearerAuth": []}],
                "consumes": ["application/json"],
                "produces": ["application/json"],
                "parameters": [
                    {"name": "id", "in": "path", "required": true, "type": "string", "format": "uuid"},
                    {
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {"$ref": "#/definitions/UpdateSystemUser"}
                    }
                ],
                "responses": {
                    "200": {"description": "Обновлено"},
                    "400": {"description": "Ошибка валидации"},
                    "403": {"description": "Нет доступа (только admin)"},
                    "404": {"description": "Пользователь не найден"},
                    "500": {"description": "Internal Server Error"}
                }
            },
            "delete": {
                "tags": ["System Users"],
                "summary": "Удалить системного пользователя",
                "security": [{"BearerAuth": []}],
                "produces": ["application/json"],
                "parameters": [
                    {"name": "id", "in": "path", "required": true, "type": "string", "format": "uuid"}
                ],
                "responses": {
                    "200": {"description": "Удалено"},
                    "400": {"description": "Пользователь не найден"},
                    "403": {"description": "Нет доступа (только admin)"},
                    "500": {"description": "Internal Server Error"}
                }
            }
        },
        "/admin/notifications": {
            "get": {
                "tags": ["Notifications"],
                "summary": "Получить историю уведомлений текущего пользователя",
                "description": "Возвращает пагинированный список системных уведомлений, адресованных авторизованному admin или supervisor.",
                "security": [{"BearerAuth": []}],
                "produces": ["application/json"],
                "parameters": [
                    {"name": "limit",  "in": "query", "type": "integer", "default": 20, "description": "Количество записей"},
                    {"name": "offset", "in": "query", "type": "integer", "default": 0,  "description": "Смещение"}
                ],
                "responses": {
                    "200": {
                        "description": "Список уведомлений",
                        "schema": {
                            "type": "object",
                            "properties": {
                                "data": {
                                    "type": "array",
                                    "items": {"$ref": "#/definitions/SystemNotification"}
                                },
                                "total":  {"type": "integer"},
                                "limit":  {"type": "integer"},
                                "offset": {"type": "integer"}
                            }
                        }
                    },
                    "401": {"description": "Не авторизован"},
                    "403": {"description": "Нет доступа"},
                    "500": {"description": "Internal Server Error"}
                }
            },
            "post": {
                "tags": ["Notifications"],
                "summary": "Отправить push-уведомление сотруднику или группе",
                "description": "Создаёт FCM push-уведомление для конкретного сотрудника или всей группы. Нужно указать ровно одно из полей: user_id или group_id.",
                "security": [{"BearerAuth": []}],
                "consumes": ["application/json"],
                "produces": ["application/json"],
                "parameters": [
                    {
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {"$ref": "#/definitions/CreateNotificationRequest"}
                    }
                ],
                "responses": {
                    "200": {"description": "Уведомление отправлено"},
                    "400": {"description": "Ошибка валидации (отсутствует user_id/group_id или оба сразу)"},
                    "401": {"description": "Не авторизован"},
                    "403": {"description": "Нет доступа"},
                    "500": {"description": "Internal Server Error"}
                }
            }
        },
        "/admin/notifications/{id}/read": {
            "patch": {
                "tags": ["Notifications"],
                "summary": "Пометить уведомление как прочитанное",
                "security": [{"BearerAuth": []}],
                "produces": ["application/json"],
                "parameters": [
                    {"name": "id", "in": "path", "required": true, "type": "string", "format": "uuid", "description": "ID уведомления"}
                ],
                "responses": {
                    "200": {"description": "Успешно помечено"},
                    "400": {"description": "Невалидный UUID"},
                    "401": {"description": "Не авторизован"},
                    "403": {"description": "Нет доступа"},
                    "404": {"description": "Уведомление не найдено"},
                    "500": {"description": "Internal Server Error"}
                }
            }
        }
    },
    "definitions": {
        "Auth": {
            "type": "object",
            "required": ["login", "password"],
            "properties": {
                "login": {"type": "string", "example": "admin"},
                "password": {"type": "string", "example": "password123"}
            }
        },
        "CreateUser": {
            "type": "object",
            "required": ["full_name", "login", "password", "confirm_password"],
            "properties": {
                "full_name": {"type": "string"},
                "login": {"type": "string"},
                "password": {"type": "string", "minLength": 8, "maxLength": 20},
                "confirm_password": {"type": "string", "minLength": 8, "maxLength": 20},
                "role": {"type": "string"},
                "group_id": {"type": "string", "format": "uuid"},
                "phone_number": {"type": "string"},
                "inn": {"type": "string"},
                "birthday": {"type": "string", "format": "date-time"}
            }
        },
        "TokenPair": {
            "type": "object",
            "properties": {
                "access_token": {"type": "string"},
                "refresh_token": {"type": "string"}
            }
        },
        "RefreshRequest": {
            "type": "object",
            "required": ["refresh_token"],
            "properties": {
                "refresh_token": {"type": "string"}
            }
        },
        "StaffUser": {
            "type": "object",
            "properties": {
                "id": {"type": "string", "format": "uuid"},
                "group_id": {"type": "string", "format": "uuid"},
                "full_name": {"type": "string"},
                "login": {"type": "string"},
                "phone_number": {"type": "string"},
                "role": {"type": "string"},
                "inn": {"type": "string"},
                "birthday": {"type": "string", "format": "date-time"},
                "is_temporary_password": {"type": "boolean"},
                "password_updated_at": {"type": "string", "format": "date-time"},
                "device_id": {"type": "string"},
                "fcm_token": {"type": "string"},
                "app_version": {"type": "string"},
                "last_battery_level": {"type": "integer"},
                "last_seen_at": {"type": "string", "format": "date-time"},
                "is_active": {"type": "boolean"},
                "deactivated_at": {"type": "string", "format": "date-time"},
                "deactivated_by_id": {"type": "string", "format": "uuid"},
                "deactivation_reason": {"type": "string"},
                "created_at": {"type": "string", "format": "date-time"}
            }
        },
        "CreateStaffUser": {
            "type": "object",
            "required": ["group_id", "full_name", "login", "password", "confirm_password"],
            "properties": {
                "group_id": {"type": "string", "format": "uuid"},
                "full_name": {"type": "string"},
                "login": {"type": "string"},
                "phone_number": {"type": "string"},
                "inn": {"type": "string"},
                "birthday": {"type": "string", "format": "date-time"},
                "password": {"type": "string"},
                "confirm_password": {"type": "string"}
            }
        },
        "UpdateStaffUser": {
            "type": "object",
            "properties": {
                "full_name": {"type": "string"},
                "phone_number": {"type": "string"},
                "is_active": {"type": "boolean"},
                "group_id": {"type": "string", "format": "uuid"},
                "fcm_token": {"type": "string"},
                "device_id": {"type": "string"},
                "deactivated_at": {"type": "string", "format": "date-time"},
                "deactivated_by_id": {"type": "string", "format": "uuid"},
                "deactivation_reason": {"type": "string"}
            }
        },
        "DeactivateStaffUser": {
            "type": "object",
            "properties": {
                "deactivated_at": {"type": "string", "format": "date-time"},
                "deactivated_by_id": {"type": "string", "format": "uuid"},
                "deactivation_reason": {"type": "string"}
            }
        },
        "StaffUserLocation": {
            "type": "object",
            "properties": {
                "id": {"type": "string", "format": "uuid"},
                "group_id": {"type": "string", "format": "uuid"},
                "full_name": {"type": "string"},
                "login": {"type": "string"},
                "phone_number": {"type": "string"},
                "inn": {"type": "string"},
                "birthday": {"type": "string", "format": "date-time"},
                "device_id": {"type": "string"},
                "fcm_token": {"type": "string"},
                "app_version": {"type": "string"},
                "last_battery_level": {"type": "integer"},
                "last_seen_at": {"type": "string", "format": "date-time"},
                "location": {"$ref": "#/definitions/LocationHistory"}
            }
        },
        "StaffUserLocationHistory": {
            "type": "object",
            "properties": {
                "id": {"type": "string", "format": "uuid"},
                "group_id": {"type": "string", "format": "uuid"},
                "full_name": {"type": "string"},
                "login": {"type": "string"},
                "phone_number": {"type": "string"},
                "inn": {"type": "string"},
                "birthday": {"type": "string", "format": "date-time"},
                "device_id": {"type": "string"},
                "fcm_token": {"type": "string"},
                "app_version": {"type": "string"},
                "last_battery_level": {"type": "integer"},
                "last_seen_at": {"type": "string", "format": "date-time"},
                "location": {
                    "type": "array",
                    "items": {"$ref": "#/definitions/LocationHistory"}
                }
            }
        },
        "LocationHistory": {
            "type": "object",
            "properties": {
                "id": {"type": "string", "format": "uuid"},
                "staff_id": {"type": "string", "format": "uuid"},
                "latitude": {"type": "number", "format": "double"},
                "longitude": {"type": "number", "format": "double"},
                "battery_level": {"type": "integer"},
                "captured_at": {"type": "string", "format": "date-time"},
                "status": {"type": "string"}
            }
        },
        "Store": {
            "type": "object",
            "properties": {
                "id": {"type": "string", "format": "uuid"},
                "name": {"type": "string"},
                "address": {"type": "string"},
                "latitude": {"type": "number", "format": "double"},
                "longitude": {"type": "number", "format": "double"},
                "owner_full_name": {"type": "string"},
                "owner_inn": {"type": "string"},
                "owner_passport": {"type": "string"},
                "owner_phone": {"type": "string"},
                "additional_phone": {"type": "string"},
                "interior_photo": {"type": "string"},
                "exterior_photo": {"type": "string"},
                "comment": {"type": "string"},
                "status": {"type": "string", "enum": ["pending", "approved", "rejected"]},
                "rejection_comment": {"type": "string"},
                "created_by_staff_id": {"type": "string", "format": "uuid"},
                "is_active": {"type": "boolean"},
                "is_frozen": {"type": "boolean"},
                "deactivated_at": {"type": "string", "format": "date-time"},
                "deactivated_by_id": {"type": "string", "format": "uuid"},
                "deactivation_reason": {"type": "string"},
                "created_at": {"type": "string", "format": "date-time"}
            }
        },
        "UpdateStore": {
            "type": "object",
            "properties": {
                "name": {"type": "string"},
                "address": {"type": "string"},
                "latitude": {"type": "number", "format": "double"},
                "longitude": {"type": "number", "format": "double"},
                "owner_full_name": {"type": "string"},
                "owner_inn": {"type": "string"},
                "owner_passport": {"type": "string"},
                "owner_phone": {"type": "string"},
                "additional_phone": {"type": "string"},
                "interior_photo": {"type": "string"},
                "exterior_photo": {"type": "string"},
                "comment": {"type": "string"},
                "status": {"type": "string", "enum": ["approved", "rejected"]},
                "rejection_comment": {"type": "string"},
                "is_active": {"type": "boolean"},
                "is_frozen": {"type": "boolean"}
            }
        },
        "StoreListFilter": {
            "type": "object",
            "properties": {
                "status": {"type": "string"},
                "is_active": {"type": "boolean"},
                "limit": {"type": "integer"},
                "offset": {"type": "integer"},
                "sort_by": {"type": "string"},
                "sort_order": {"type": "string", "enum": ["asc", "desc"]}
            }
        },
        "Group": {
            "type": "object",
            "properties": {
                "id": {"type": "string", "format": "uuid"},
                "name": {"type": "string"},
                "supervisor_id": {"type": "string", "format": "uuid"},
                "supervisor_name": {"type": "string"},
                "created_at": {"type": "string", "format": "date-time"}
            }
        },
        "Visit": {
            "type": "object",
            "properties": {
                "id": {"type": "string", "format": "uuid"},
                "schedule_id": {"type": "string", "format": "uuid"},
                "staff_id": {"type": "string", "format": "uuid"},
                "store_id": {"type": "string", "format": "uuid"},
                "check_in_at": {"type": "string", "format": "date-time"},
                "check_out_at": {"type": "string", "format": "date-time"},
                "check_in_lat": {"type": "number", "format": "double"},
                "check_in_lon": {"type": "number", "format": "double"},
                "distance_deviation": {"type": "integer"},
                "selfie_url": {"type": "string"},
                "comment": {"type": "string"},
                "is_offline": {"type": "boolean"},
                "review_status": {"type": "string", "enum": ["new", "pending", "approved", "rejected", "revision"]},
                "review_comment": {"type": "string"},
                "reviewer_id": {"type": "string", "format": "uuid"},
                "reviewed_at": {"type": "string", "format": "date-time"},
                "created_at": {"type": "string", "format": "date-time"}
            }
        },
        "ReviewVisitRequest": {
            "type": "object",
            "required": ["status"],
            "properties": {
                "status": {"type": "string", "enum": ["approved", "rejected", "revision"], "description": "approved и rejected — финальные (неизменяемые). revision — на доработку, сотрудник может дополнить визит."},
                "review_comment": {"type": "string"}
            }
        },
        "CreateSchedule": {
            "type": "object",
            "required": ["staff_id", "store_id", "planned_date", "description", "window_open", "window_close", "created_by_admin_id"],
            "properties": {
                "staff_id": {"type": "string", "format": "uuid"},
                "store_id": {"type": "string", "format": "uuid"},
                "planned_date": {"type": "string", "format": "date-time"},
                "description": {"type": "string"},
                "window_open": {"type": "string", "format": "date-time"},
                "window_close": {"type": "string", "format": "date-time"},
                "is_photo_required": {"type": "boolean"},
                "is_selfie_required": {"type": "boolean"},
                "is_comment_required": {"type": "boolean"},
                "created_by_admin_id": {"type": "string", "format": "uuid"}
            }
        },
        "SystemSettings": {
            "type": "object",
            "properties": {
                "id": {"type": "string", "format": "uuid"},
                "key": {"type": "string"},
                "value": {"type": "string"},
                "description": {"type": "string"}
            }
        },
        "UpdateSystemSettings": {
            "type": "object",
            "properties": {
                "key": {"type": "string"},
                "value": {"type": "string"},
                "description": {"type": "string"}
            }
        },
        "AuditLog": {
            "type": "object",
            "properties": {
                "id": {"type": "string", "format": "uuid"},
                "actor_id": {"type": "string", "format": "uuid"},
                "actor_type": {"type": "string"},
                "action": {"type": "string"},
                "target_table": {"type": "string"},
                "target_id": {"type": "string", "format": "uuid"},
                "old_value": {"type": "object"},
                "new_value": {"type": "object"},
                "ip_address": {"type": "string"},
                "created_at": {"type": "string", "format": "date-time"}
            }
        },
        "SystemUser": {
            "type": "object",
            "properties": {
                "id": {"type": "string", "format": "uuid"},
                "full_name": {"type": "string"},
                "login": {"type": "string"},
                "role": {"type": "string", "enum": ["admin", "supervisor"]},
                "is_active": {"type": "boolean"},
                "deactivated_at": {"type": "string", "format": "date-time"},
                "deactivated_by_id": {"type": "integer"},
                "deactivation_reason": {"type": "string"},
                "created_at": {"type": "string", "format": "date-time"}
            }
        },
        "CreateSystemUser": {
            "type": "object",
            "required": ["full_name", "login", "password", "confirm_password", "role"],
            "properties": {
                "full_name": {"type": "string"},
                "login": {"type": "string"},
                "password": {"type": "string", "minLength": 8, "maxLength": 20},
                "confirm_password": {"type": "string", "minLength": 8, "maxLength": 20},
                "role": {"type": "string", "enum": ["admin", "supervisor"]}
            }
        },
        "UpdateSystemUser": {
            "type": "object",
            "properties": {
                "full_name": {"type": "string"},
                "is_active": {"type": "boolean"},
                "deactivated_at": {"type": "string", "format": "date-time"},
                "deactivated_by_id": {"type": "integer"},
                "deactivation_reason": {"type": "string"}
            }
        },
        "StaffUserSchedule": {
            "type": "object",
            "properties": {
                "id": {"type": "string", "format": "uuid"},
                "staff_id": {"type": "string", "format": "uuid"},
                "weekday": {"type": "integer", "description": "День недели (0 = воскресенье, 1 = понедельник, ..., 6 = суббота)"},
                "start_time": {"type": "string", "format": "date-time"},
                "end_time": {"type": "string", "format": "date-time"},
                "created_at": {"type": "string", "format": "date-time"}
            }
        },
        "UpdateStaffSchedule": {
            "type": "object",
            "properties": {
                "weekday": {"type": "integer", "description": "День недели (0-6)"},
                "start_time": {"type": "string", "format": "date-time"},
                "end_time": {"type": "string", "format": "date-time"}
            }
        },
        "Schedule": {
            "type": "object",
            "properties": {
                "id": {"type": "string", "format": "uuid"},
                "staff_id": {"type": "string", "format": "uuid"},
                "store_id": {"type": "string", "format": "uuid"},
                "planned_date": {"type": "string", "format": "date-time"},
                "description": {"type": "string"},
                "window_open": {"type": "string", "format": "date-time"},
                "window_close": {"type": "string", "format": "date-time"},
                "is_photo_required": {"type": "boolean"},
                "is_selfie_required": {"type": "boolean"},
                "is_comment_required": {"type": "boolean"},
                "status": {"type": "string"},
                "created_by_admin_id": {"type": "string", "format": "uuid"}
            }
        },
        "VisitWithPhoto": {
            "type": "object",
            "properties": {
                "visit": {"$ref": "#/definitions/Visit"},
                "photos": {
                    "type": "array",
                    "items": {"$ref": "#/definitions/VisitPhotos"}
                }
            }
        },
        "VisitPhotos": {
            "type": "object",
            "properties": {
                "id": {"type": "string", "format": "uuid"},
                "visit_id": {"type": "string", "format": "uuid"},
                "photo_url": {"type": "string"},
                "captured_at": {"type": "string", "format": "date-time"},
                "original_name": {"type": "string"},
                "photo_type": {"type": "string", "enum": ["selfie", "shelf"]}
            }
        },
        "Pagination": {
            "type": "object",
            "properties": {
                "limit": {"type": "integer"},
                "offset": {"type": "integer"},
                "total": {"type": "integer"}
            }
        },
        "ErrorResponse": {
            "type": "object",
            "properties": {
                "error": {"type": "string"},
                "status_code": {"type": "integer"},
                "request_id": {"type": "string"}
            }
        },
        "SystemNotification": {
            "type": "object",
            "properties": {
                "id":           {"type": "string", "format": "uuid"},
                "recipient_id": {"type": "string", "format": "uuid"},
                "staff_id":     {"type": "string", "format": "uuid"},
                "group_id":     {"type": "string", "format": "uuid"},
                "type":         {"type": "string", "example": "gps_missing"},
                "title":        {"type": "string"},
                "body":         {"type": "string"},
                "is_read":      {"type": "boolean"},
                "read_at":      {"type": "string", "format": "date-time"},
                "created_at":   {"type": "string", "format": "date-time"}
            }
        },
        "CreateNotificationRequest": {
            "type": "object",
            "required": ["title", "body", "type"],
            "properties": {
                "user_id":  {"type": "string", "format": "uuid", "description": "ID сотрудника (исключает group_id)"},
                "group_id": {"type": "string", "format": "uuid", "description": "ID группы (исключает user_id)"},
                "title":    {"type": "string"},
                "body":     {"type": "string"},
                "type":     {"type": "string"},
                "meta":     {"type": "object", "additionalProperties": true}
            }
        }
    }
}`

type s struct{}

func (s *s) ReadDoc() string {
	return docTemplate
}

func init() {
	swag.Register("admin", &s{})
}

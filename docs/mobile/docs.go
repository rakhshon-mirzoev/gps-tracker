package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "swagger": "2.0",
    "info": {
        "title": "GPS Tracker — Mobile API",
        "description": "Мобильное приложение GPS-трекера: авторизация сотрудников, профиль, расписания, визиты, геолокация, уведомления.",
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
        {"name": "Auth", "description": "Аутентификация сотрудников"},
        {"name": "Profile", "description": "Профиль сотрудника"},
        {"name": "Notifications", "description": "Push-уведомления"},
        {"name": "Schedule", "description": "Расписание визитов"},
        {"name": "Stores", "description": "Торговые точки"},
        {"name": "Location", "description": "Отправка геолокации"},
        {"name": "Visits", "description": "Визиты в торговые точки"},
        {"name": "Settings", "description": "Системные настройки"},
        {"name": "Work Tracking", "description": "Начало и окончание рабочего дня"}
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
        "/settings": {
            "get": {
                "tags": ["Settings"],
                "summary": "Список системных настроек (публичный)",
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
        "/auth/staff/login": {
            "post": {
                "tags": ["Auth"],
                "summary": "Вход в систему (сотрудник)",
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
        "/auth/staff/register": {
            "post": {
                "tags": ["Auth"],
                "summary": "Регистрация сотрудника",
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
        "/auth/staff/refresh": {
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
        "/auth/staff/logout": {
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
        "/mobile/profile": {
            "get": {
                "tags": ["Profile"],
                "summary": "Получить профиль текущего сотрудника",
                "security": [{"BearerAuth": []}],
                "produces": ["application/json"],
                "parameters": [
                    {"name": "limit", "in": "query", "type": "integer", "default": 20},
                    {"name": "offset", "in": "query", "type": "integer", "default": 0}
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "object",
                            "properties": {
                                "profile": {"$ref": "#/definitions/StaffUserProfile"}
                            }
                        }
                    },
                    "403": {"description": "Нет доступа"},
                    "500": {"description": "Internal Server Error"}
                }
            },
            "put": {
                "tags": ["Profile"],
                "summary": "Обновить профиль сотрудника",
                "security": [{"BearerAuth": []}],
                "consumes": ["application/json"],
                "produces": ["application/json"],
                "parameters": [
                    {
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {"$ref": "#/definitions/UpdateProfile"}
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Обновлено",
                        "schema": {
                            "type": "object",
                            "properties": {
                                "message": {"type": "string", "example": "updated"}
                            }
                        }
                    },
                    "400": {"description": "Ошибка валидации"},
                    "403": {"description": "Нет доступа"},
                    "500": {"description": "Internal Server Error"}
                }
            }
        },
        "/mobile/profile/avatar": {
            "put": {
                "tags": ["Profile"],
                "summary": "Обновить аватар сотрудника",
                "security": [{"BearerAuth": []}],
                "consumes": ["multipart/form-data"],
                "produces": ["application/json"],
                "parameters": [
                    {"name": "avatar", "in": "formData", "required": true, "type": "file", "description": "Файл изображения (до 5 МБ)"}
                ],
                "responses": {
                    "200": {
                        "description": "Аватар обновлён",
                        "schema": {
                            "type": "object",
                            "properties": {
                                "message": {"type": "string", "example": "avatar updated"}
                            }
                        }
                    },
                    "400": {"description": "Ошибка загрузки файла"},
                    "500": {"description": "Internal Server Error"}
                }
            }
        },
        "/mobile/device-token": {
            "put": {
                "tags": ["Profile"],
                "summary": "Обновить FCM / device token",
                "security": [{"BearerAuth": []}],
                "consumes": ["application/json"],
                "produces": ["application/json"],
                "parameters": [
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
                    "404": {"description": "Пользователь не найден"},
                    "500": {"description": "Internal Server Error"}
                }
            }
        },
        "/mobile/notifications": {
            "get": {
                "tags": ["Notifications"],
                "summary": "Список уведомлений текущего пользователя",
                "security": [{"BearerAuth": []}],
                "produces": ["application/json"],
                "parameters": [
                    {"name": "limit", "in": "query", "required": true, "type": "integer"},
                    {"name": "offset", "in": "query", "required": true, "type": "integer"}
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "object",
                            "properties": {
                                "notifications": {
                                    "type": "array",
                                    "items": {"$ref": "#/definitions/Notification"}
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
        "/mobile/notifications/{id}/read": {
            "patch": {
                "tags": ["Notifications"],
                "summary": "Отметить уведомление как прочитанное",
                "security": [{"BearerAuth": []}],
                "produces": ["application/json"],
                "parameters": [
                    {"name": "id", "in": "path", "required": true, "type": "string", "format": "uuid"}
                ],
                "responses": {
                    "200": {"description": "OK"},
                    "403": {"description": "Нет доступа"},
                    "404": {"description": "Уведомление не найдено"},
                    "500": {"description": "Internal Server Error"}
                }
            }
        },
        "/mobile/schedule": {
            "get": {
                "tags": ["Schedule"],
                "summary": "Расписание текущего сотрудника",
                "security": [{"BearerAuth": []}],
                "produces": ["application/json"],
                "parameters": [
                    {"name": "limit", "in": "query", "required": true, "type": "integer"},
                    {"name": "offset", "in": "query", "required": true, "type": "integer"}
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "object",
                            "properties": {
                                "schedules": {
                                    "type": "array",
                                    "items": {"$ref": "#/definitions/Schedule"}
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
        "/mobile/my-stores": {
            "get": {
                "tags": ["Stores"],
                "summary": "Магазины, привязанные к текущему сотруднику",
                "security": [{"BearerAuth": []}],
                "produces": ["application/json"],
                "parameters": [
                    {"name": "limit", "in": "query", "required": true, "type": "integer"},
                    {"name": "offset", "in": "query", "required": true, "type": "integer"}
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
                                }
                            }
                        }
                    },
                    "403": {"description": "Нет доступа"},
                    "500": {"description": "Internal Server Error"}
                }
            }
        },
        "/mobile/start-work": {
            "put": {
                "tags": ["Work Tracking"],
                "summary": "Начать рабочий день",
                "security": [{"BearerAuth": []}],
                "produces": ["application/json"],
                "responses": {
                    "204": {"description": "Рабочий день начат"},
                    "400": {"description": "Ошибка"},
                    "500": {"description": "Internal Server Error"}
                }
            }
        },
        "/mobile/end-work": {
            "put": {
                "tags": ["Work Tracking"],
                "summary": "Завершить рабочий день",
                "security": [{"BearerAuth": []}],
                "produces": ["application/json"],
                "responses": {
                    "204": {"description": "Рабочий день завершён"},
                    "400": {"description": "Ошибка"},
                    "500": {"description": "Internal Server Error"}
                }
            }
        },
        "/mobile/stores/upload-photo": {
            "post": {
                "tags": ["Stores"],
                "summary": "Загрузить фото магазина",
                "security": [{"BearerAuth": []}],
                "consumes": ["multipart/form-data"],
                "produces": ["application/json"],
                "parameters": [
                    {"name": "photo", "in": "formData", "required": true, "type": "file", "description": "Фото магазина (png/jpg/jpeg, до 10 МБ)"}
                ],
                "responses": {
                    "200": {
                        "description": "Фото загружено",
                        "schema": {
                            "type": "object",
                            "properties": {
                                "photo_key": {"type": "string"},
                                "download_url": {"type": "string"}
                            }
                        }
                    },
                    "400": {"description": "Ошибка загрузки файла"},
                    "500": {"description": "Internal Server Error"}
                }
            }
        },
        "/mobile/stores/propose": {
            "post": {
                "tags": ["Stores"],
                "summary": "Предложить новый магазин",
                "security": [{"BearerAuth": []}],
                "consumes": ["application/json"],
                "produces": ["application/json"],
                "parameters": [
                    {
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {"$ref": "#/definitions/CreateStore"}
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Магазин создан",
                        "schema": {
                            "type": "object",
                            "properties": {
                                "message": {"type": "string", "example": "create complete"}
                            }
                        }
                    },
                    "400": {"description": "Ошибка валидации"},
                    "500": {"description": "Internal Server Error"}
                }
            }
        },
        "/mobile/location": {
            "post": {
                "tags": ["Location"],
                "summary": "Отправить пакет геолокаций",
                "security": [{"BearerAuth": []}],
                "consumes": ["application/json"],
                "produces": ["application/json"],
                "parameters": [
                    {
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "type": "array",
                            "items": {"$ref": "#/definitions/LocationHistory"}
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Локации сохранены",
                        "schema": {
                            "type": "object",
                            "properties": {
                                "message": {"type": "string", "example": "locations created"},
                                "count": {"type": "integer"}
                            }
                        }
                    },
                    "400": {"description": "Пустой массив или ошибка валидации"},
                    "500": {"description": "Internal Server Error"}
                }
            }
        },
        "/mobile/visit": {
            "get": {
                "tags": ["Visits"],
                "summary": "Список визитов текущего сотрудника",
                "security": [{"BearerAuth": []}],
                "produces": ["application/json"],
                "parameters": [
                    {"name": "limit", "in": "query", "type": "integer", "default": 20},
                    {"name": "offset", "in": "query", "type": "integer", "default": 0}
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "object",
                            "properties": {
                                "data": {
                                    "type": "array",
                                    "items": {"$ref": "#/definitions/Visit"}
                                },
                                "pagination": {"$ref": "#/definitions/Pagination"}
                            }
                        }
                    },
                    "403": {"description": "Нет доступа"},
                    "500": {"description": "Internal Server Error"}
                }
            }
        },
        "/mobile/visit/check-in": {
            "post": {
                "tags": ["Visits"],
                "summary": "Начать визит (check-in)",
                "security": [{"BearerAuth": []}],
                "consumes": ["application/json"],
                "produces": ["application/json"],
                "parameters": [
                    {
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {"$ref": "#/definitions/CreateVisit"}
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Check-in выполнен",
                        "schema": {
                            "type": "object",
                            "properties": {
                                "message": {"type": "string", "example": "check-in completed"}
                            }
                        }
                    },
                    "400": {"description": "Уже рассмотрен или ошибка валидации"},
                    "500": {"description": "Internal Server Error"}
                }
            }
        },
        "/mobile/visit/{id}/photo": {
            "post": {
                "tags": ["Visits"],
                "summary": "Загрузить фото визита",
                "security": [{"BearerAuth": []}],
                "consumes": ["multipart/form-data"],
                "produces": ["application/json"],
                "parameters": [
                    {"name": "id", "in": "path", "required": true, "type": "string", "description": "ID визита"}
                ],
                "responses": {
                    "200": {
                        "description": "Фото загружено",
                        "schema": {
                            "type": "object",
                            "properties": {
                                "message": {"type": "string", "example": "update complete"}
                            }
                        }
                    },
                    "400": {"description": "Ошибка валидации"},
                    "500": {"description": "Internal Server Error"}
                }
            }
        },
        "/mobile/visit/{id}/complete": {
            "post": {
                "tags": ["Visits"],
                "summary": "Завершить визит (check-out)",
                "security": [{"BearerAuth": []}],
                "consumes": ["application/json"],
                "produces": ["application/json"],
                "parameters": [
                    {"name": "id", "in": "path", "required": true, "type": "string", "description": "ID визита"},
                    {
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {"$ref": "#/definitions/UpdateVisit"}
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Визит завершён",
                        "schema": {
                            "type": "object",
                            "properties": {
                                "message": {"type": "string", "example": "visit complete"}
                            }
                        }
                    },
                    "400": {"description": "Ошибка валидации"},
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
                "login": {"type": "string", "example": "staff_user"},
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
        "StaffUserProfile": {
            "type": "object",
            "properties": {
                "user": {"$ref": "#/definitions/StaffUser"},
                "visits": {
                    "type": "array",
                    "items": {"$ref": "#/definitions/Visit"}
                },
                "total_count": {"type": "integer"}
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
        "UpdateStaffUser": {
            "type": "object",
            "properties": {
                "fcm_token": {"type": "string"},
                "device_id": {"type": "string"}
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
        "CreateStore": {
            "type": "object",
            "required": ["name", "address", "latitude", "longitude", "owner_full_name", "owner_inn", "owner_phone"],
            "properties": {
                "name": {"type": "string"},
                "address": {"type": "string"},
                "latitude": {"type": "number", "format": "double"},
                "longitude": {"type": "number", "format": "double"},
                "owner_full_name": {"type": "string"},
                "owner_inn": {"type": "string"},
                "owner_phone": {"type": "string"},
                "owner_passport": {"type": "string"},
                "additional_phone": {"type": "string"},
                "interior_photo": {"type": "string"},
                "exterior_photo": {"type": "string"},
                "comment": {"type": "string"}
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
        "CreateVisit": {
            "type": "object",
            "required": ["schedule_id", "store_id", "lat", "lon"],
            "properties": {
                "schedule_id": {"type": "string", "format": "uuid"},
                "store_id": {"type": "string", "format": "uuid"},
                "lat": {"type": "number", "format": "double"},
                "lon": {"type": "number", "format": "double"}
            }
        },
        "UpdateVisit": {
            "type": "object",
            "properties": {
                "check_out_at": {"type": "string", "format": "date-time"},
                "selfie_url": {"type": "string"},
                "comment": {"type": "string"}
            }
        },
        "Notification": {
            "type": "object",
            "properties": {
                "id": {"type": "string", "format": "uuid"},
                "user_id": {"type": "string", "format": "uuid"},
                "type": {"type": "string"},
                "title": {"type": "string"},
                "body": {"type": "string"},
                "meta": {"type": "object"},
                "is_read": {"type": "boolean"},
                "push_sent": {"type": "boolean"},
                "read_at": {"type": "string", "format": "date-time"},
                "push_sent_at": {"type": "string", "format": "date-time"},
                "created_at": {"type": "string", "format": "date-time"}
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
        "UpdateProfile": {
            "type": "object",
            "properties": {
                "full_name": {"type": "string"},
                "phone_number": {"type": "string"}
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
                "review_status": {"type": "string", "enum": ["pending", "approved", "rejected"]},
                "review_comment": {"type": "string"},
                "reviewer_id": {"type": "string", "format": "uuid"},
                "reviewed_at": {"type": "string", "format": "date-time"},
                "created_at": {"type": "string", "format": "date-time"}
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
        }
    }
}`

type s struct{}

func (s *s) ReadDoc() string {
	return docTemplate
}

func init() {
	swag.Register("mobile", &s{})
}


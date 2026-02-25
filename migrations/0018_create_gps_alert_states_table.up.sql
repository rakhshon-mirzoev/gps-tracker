-- Хранит состояние GPS-алертов: один активный алерт на сотрудника.
-- Запись создаётся при отправке уведомления, удаляется когда сотрудник
-- снова начинает присылать GPS (в location_history.Create).
CREATE TABLE IF NOT EXISTS gps_alert_states (
    staff_id    UUID PRIMARY KEY
        REFERENCES staff_users(id)
        ON DELETE CASCADE,
    alerted_at  TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now()
);

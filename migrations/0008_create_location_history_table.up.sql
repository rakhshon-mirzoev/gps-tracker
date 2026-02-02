CREATE TABLE location_history(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    staff_id UUID
        REFERENCES staff_users(id)
        ON DELETE RESTRICT,
    latitude DECIMAL(9,6)  NOT NULL,
    longitude DECIMAL(9,6) NOT NULL,
    battery_level INTEGER NOT NULL,
    captured_at TIMESTAMP NOT NULL DEFAULT now(),
    status TEXT NOT NULL,

    CONSTRAINT location_history_status_check
        CHECK (status IN ('online','offline'))
);

CREATE INDEX idx_location_history_staff_time
    ON location_history (staff_id, captured_at);

CREATE INDEX idx_location_history_staff_status_time
    ON location_history (staff_id, status, captured_at);

CREATE INDEX idx_location_history_captured_at
    ON location_history (captured_at);
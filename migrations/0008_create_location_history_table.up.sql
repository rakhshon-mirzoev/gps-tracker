CREATE TABLE location_history(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    staff_id UUID
        REFERENCES staff_users(id)
        ON DELETE RESTRICT,
    latitude DECIMAL(9,6)  NOT NULL,
    longitude DECIMAL(9,6) NOT NULL,
    battery_level INTEGER NOT NULL,
    captured_at TIMESTAMP NOT NULL DEFAULT now(),
    status TEXT NOT NULL
);
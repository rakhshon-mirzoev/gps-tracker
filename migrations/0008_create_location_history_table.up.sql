CREATE TABLE location_history(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    staff_id UUID
        REFERENCES staff_users(id)
        ON DELETE RESTRICT,
    location geography(POINT, 4326)   NOT NULL,
    battery_level INTEGER,
    captured_at TIMESTAMP
);

CREATE INDEX idx_location_history_gist 
    ON location_history USING GIST (location);
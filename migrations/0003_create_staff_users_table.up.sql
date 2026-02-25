CREATE TABLE staff_users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    group_id UUID
        REFERENCES groups(id)
        ON DELETE RESTRICT,

    full_name VARCHAR(255) NOT NULL,
    login VARCHAR(100) NOT NULL UNIQUE,
    phone_number VARCHAR(50),
    inn VARCHAR(20),
    role TEXT NOT NULL DEFAULT 'staff',

    birthday DATE,

    password_hash TEXT NOT NULL,
    is_temporary_password BOOLEAN NOT NULL DEFAULT false,
    password_updated_at TIMESTAMP,

    device_id VARCHAR(255),
    fcm_token VARCHAR(255),
    app_version VARCHAR(50),

    last_battery_level INTEGER
        CHECK (last_battery_level BETWEEN 0 AND 100),

    last_seen_at TIMESTAMP,

    is_active BOOLEAN NOT NULL DEFAULT true,

    deactivated_at TIMESTAMP,
    deactivated_by_id UUID
        REFERENCES system_users(id)
        ON DELETE SET NULL,

    deactivation_reason TEXT,

    avatar_key TEXT,
    start_at TIMESTAMP,
    end_at TIMESTAMP,

    created_at TIMESTAMP NOT NULL DEFAULT now()
);

CREATE INDEX idx_location_history_staff_captured_at
ON location_history (staff_id, captured_at DESC);

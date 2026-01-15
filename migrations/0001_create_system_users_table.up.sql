CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE system_users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    full_name TEXT NOT NULL ,
    password_hash TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT now(),
    login TEXT NOT NULL,
    role TEXT NOT NULL,
    is_active BOOLEAN DEFAULT false,
    deactivated_at TIMESTAMP DEFAULT NULL,
    deactivated_by_id UUID DEFAULT NULL,
    deactivation_reason TEXT DEFAULT NULL,

    CONSTRAINT fk_deactivated_by
        FOREIGN KEY (deactivated_by_id)
        REFERENCES system_users(id)
        ON DELETE SET NULL
);
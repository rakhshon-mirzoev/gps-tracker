CREATE TABLE system_users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    full_name TEXT NOT NULL ,
    password_hash TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT now(),
    login TEXT NOT NULL,
    role TEXT NOT NULL,
    is_active BOOLEAN DEFAULT false,
    deactivated_at TIMESTAMP DEFAULT NULL,
    deactivated_by_id BIGINT DEFAULT NULL,
    deactivation_reason TEXT DEFAULT NULL
);
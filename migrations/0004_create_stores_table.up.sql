CREATE TABLE stores
(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    name VARCHAR(200)        NOT NULL,
    address TEXT                NOT NULL,

    latitude DECIMAL(9,6)   NOT NULL,
    longitude DECIMAL(9,6) NOT NULL,

    owner_full_name VARCHAR(160)        NOT NULL,
    owner_inn     CHAR(12)            NOT NULL,
    owner_passport VARCHAR(25),
    owner_phone VARCHAR(15),
    additional_phone VARCHAR(15),
    interior_photo TEXT,
    exterior_photo TEXT,
    comment TEXT,

    group_id UUID
        REFERENCES groups(id)
        ON DELETE RESTRICT,

    status   VARCHAR(40)         NOT NULL DEFAULT 'pending',
    rejection_comment  TEXT,

    created_by_staff_id UUID NOT NULL 
        REFERENCES staff_users(id)
        ON DELETE RESTRICT,

    is_active  BOOLEAN             NOT NULL DEFAULT FALSE,
    is_frozen    BOOLEAN             NOT NULL DEFAULT FALSE,

    deactivated_at    TIMESTAMP,
    deactivated_by_id  UUID 
        REFERENCES system_users(id)
        ON DELETE RESTRICT,
    deactivation_reason TEXT,

    created_at  TIMESTAMP         NOT NULL DEFAULT now(),

    CONSTRAINT stores_status_check    
        CHECK (status IN ('pending', 'approved', 'rejected', 'deactivated', 'frozen'))
);

CREATE INDEX idx_stores_status_active 
    ON stores (status, is_active)
    WHERE is_active = true;

CREATE INDEX idx_stores_owner_inn 
    ON stores (owner_inn);
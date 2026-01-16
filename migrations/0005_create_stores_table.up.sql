CREATE EXTENSION IF NOT EXISTS postgis;

CREATE TABLE stores
(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    name VARCHAR(200)        NOT NULL,
    address TEXT                NOT NULL,

    location geography(POINT, 4326)   NOT NULL,

    owner_full_name VARCHAR(160)        NOT NULL,
    owner_inn     CHAR(12)            NOT NULL,
    owner_passport VARCHAR(25),

    status   VARCHAR(40)         NOT NULL DEFAULT 'new',
    rejection_comment  TEXT,

    created_by_staff_id UUID NOT NULL 
        REFERENCES staff_users(id)
        ON DELETE RESTRICT,
    created_by_admin_id UUID NOT NULL 
        REFERENCES system_users(id)
        ON DELETE RESTRICT,

    is_active  BOOLEAN             NOT NULL DEFAULT TRUE,
    is_frozen    BOOLEAN             NOT NULL DEFAULT FALSE,

    deactivated_at    TIMESTAMPTZ,
    deactivated_by_id  UUID NOT NULL 
        REFERENCES system_users(id)
        ON DELETE RESTRICT,
    deactivation_reason TEXT,

    created_at  TIMESTAMPTZ         NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT stores_latitude_check  
        CHECK (ST_Y(location::geometry) BETWEEN -90 AND 90),
        
    CONSTRAINT stores_longitude_check 
        CHECK (ST_X(location::geometry) BETWEEN -180 AND 180),

    CONSTRAINT stores_status_check    
        CHECK (status IN ('new', 'pending', 'approved', 'rejected', 'deactivated'))
);

CREATE INDEX idx_stores_location_gist 
    ON stores USING GIST (location);

CREATE INDEX idx_stores_status_active 
    ON stores (status, is_active)
    WHERE is_active = true;

CREATE INDEX idx_stores_owner_inn 
    ON stores (owner_inn);
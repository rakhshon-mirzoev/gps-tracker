CREATE TABLE groups (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL UNIQUE ,
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    supervisor_id UUID,
    parent_id UUID REFERENCES groups(id) ON DELETE RESTRICT,
    is_deleted BOOLEAN NOT NULL DEFAULT FALSE,
    deleted_at TIMESTAMP,
    type VARCHAR(20) DEFAULT 'staff',
    color VARCHAR(7) DEFAULT '#000000',
    cluster VARCHAR(100) DEFAULT '',
    layer INTEGER DEFAULT 0,


    CONSTRAINT fk_groups_supervisor
        FOREIGN KEY (supervisor_id)
        REFERENCES system_users(id)
        ON DELETE SET NULL
);

CREATE INDEX idx_groups_supervisor_id ON groups(supervisor_id);
CREATE INDEX idx_groups_parent_id ON groups(parent_id);

ALTER TABLE groups ADD CONSTRAINT groups_type_check CHECK (type IN ('staff', 'store'));
CREATE INDEX idx_groups_type ON groups(type);
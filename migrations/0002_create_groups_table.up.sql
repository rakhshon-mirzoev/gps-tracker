CREATE TABLE groups (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL UNIQUE ,
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    supervisor_id UUID,

    CONSTRAINT fk_groups_supervisor
        FOREIGN KEY (supervisor_id)
        REFERENCES system_users(id)
        ON DELETE SET NULL
);

CREATE INDEX idx_groups_supervisor_id ON groups(supervisor_id);
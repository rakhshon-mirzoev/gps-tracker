CREATE TABLE schedules (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    staff_id UUID NOT NULL
        REFERENCES staff_users(id)
        ON DELETE RESTRICT,

    store_id UUID NOT NULL
        REFERENCES stores(id)
        ON DELETE RESTRICT,

    planned_date TIMESTAMP,
    description TEXT,
    window_open TIMESTAMP,
    window_close TIMESTAMP,

    is_photo_required BOOLEAN,
    is_selfie_required BOOLEAN,
    is_comment_required BOOLEAN,
    status varchar(30),

    created_by_admin_id UUID
        REFERENCES system_users(id)
        ON DELETE RESTRICT
);

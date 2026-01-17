CREATE TABLE daily_reporting_stats (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    date TIMESTAMP NOT NULL,
    staff_id UUID NOT NULL
        REFERENCES staff_users(id)
        ON DELETE RESTRICT,
    group_id UUID NOT NULL
        REFERENCES groups(id)
        ON DELETE RESTRICT,
    planned_count INTEGER NOT NULL,
    completed_count INTEGER NOT NULL,
    rejected_count INTEGER NOT NULL,
    missed_count INTEGER NOT NULL,

    work_start_at TIMESTAMPTZ,
    work_end_at TIMESTAMPTZ,

    total_km FLOAT
);
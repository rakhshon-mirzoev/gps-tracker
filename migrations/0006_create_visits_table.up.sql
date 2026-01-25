CREATE TABLE visits(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    schedule_id UUID
        REFERENCES schedules(id)
        ON DELETE SET NULL,
    staff_id UUID
        REFERENCES staff_users(id)
        ON DELETE RESTRICT,
    store_id UUID
        REFERENCES stores(id)
        ON DELETE RESTRICT,
    check_in_at TIMESTAMP,
    check_out_at TIMESTAMP,
    check_in_lat DECIMAL(9,6),
    check_in_lon DECIMAL(9,6),

    distance_deviation  INTEGER   DEFAULT 0,
    selfie_url   VARCHAR(512),
    comment   TEXT,
    is_offline      BOOLEAN   NOT NULL DEFAULT FALSE,

    review_status  VARCHAR(40)  NOT NULL DEFAULT 'pending',

    review_comment  TEXT,
    reviewer_id     UUID,
    reviewed_at   TIMESTAMPTZ,

    created_at  TIMESTAMPTZ  NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT visits_review_status_check
        CHECK (review_status IN ('pending', 'approved', 'rejected', 'revision_needed')),

    CONSTRAINT visits_checkout_after_checkin
        CHECK (check_out_at IS NULL OR check_out_at > check_in_at)
);

CREATE INDEX idx_visits_schedule_id     ON visits(schedule_id);
CREATE INDEX idx_visits_staff_id        ON visits(staff_id);
CREATE INDEX idx_visits_store_id        ON visits(store_id);
CREATE INDEX idx_visits_check_in_at     ON visits(check_in_at);
CREATE INDEX idx_visits_review_status   ON visits(review_status);

CREATE INDEX idx_visits_check_in_location_gist 
    ON visits USING GIST (check_in_location);

CREATE INDEX idx_visits_staff_date 
    ON visits (staff_id, check_in_at);

CREATE INDEX idx_visits_store_date_status 
    ON visits (store_id, check_in_at, review_status);
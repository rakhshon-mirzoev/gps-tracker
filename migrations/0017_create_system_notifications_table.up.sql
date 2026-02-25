CREATE TABLE IF NOT EXISTS system_notifications (
    id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    recipient_id UUID NOT NULL
        REFERENCES system_users(id)
        ON DELETE CASCADE,
    staff_id     UUID NULL
        REFERENCES staff_users(id)
        ON DELETE SET NULL,
    group_id     UUID NULL,
    type         VARCHAR(64) NOT NULL,
    title        TEXT        NOT NULL,
    body         TEXT        NOT NULL,
    is_read      BOOLEAN     NOT NULL DEFAULT FALSE,
    read_at      TIMESTAMP   NULL,
    created_at   TIMESTAMP WITH TIME ZONE DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_sys_notif_recipient_id ON system_notifications(recipient_id);
CREATE INDEX IF NOT EXISTS idx_sys_notif_created_at  ON system_notifications(created_at DESC);
CREATE INDEX IF NOT EXISTS idx_sys_notif_unread
    ON system_notifications(recipient_id)
    WHERE is_read = FALSE;

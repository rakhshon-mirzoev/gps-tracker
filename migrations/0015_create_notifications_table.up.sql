CREATE TABLE IF NOT EXISTS notifications (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NULL
        REFERENCES staff_users(id)
        ON DELETE CASCADE,
    type VARCHAR(64) NOT NULL,
    title TEXT,
    body TEXT,
    meta JSONB NULL,
    is_read BOOLEAN DEFAULT FALSE,
    push_sent BOOLEAN DEFAULT FALSE,
    read_at TIMESTAMP NULL,
    push_sent_at TIMESTAMP NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_notifications_user_id ON notifications(user_id);
CREATE INDEX IF NOT EXISTS idx_notifications_created_at ON notifications(created_at DESC);
CREATE INDEX idx_notifications_unsent
    ON notifications (push_sent)
    WHERE push_sent = FALSE;
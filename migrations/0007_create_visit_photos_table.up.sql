CREATE TABLE visit_photos (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    visit_id UUID
        REFERENCES visits(id)
        ON DELETE CASCADE,
    photo_key TEXT NOT NULL,
    original_name TEXT NOT NULL,
    captured_at TIMESTAMP NOT NULL DEFAULT now()
);
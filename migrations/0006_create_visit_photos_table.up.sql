CREATE TABLE visit_photos (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    visit_id UUID
        REFERENCES visits(id)
        ON DELETE RESTRICT,
    photo_url TEXT,
    captured_at TIMESTAMP
);
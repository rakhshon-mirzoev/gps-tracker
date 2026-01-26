CREATE TABLE visit_photos (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    visit_id UUID
        REFERENCES visits(id)
        ON DELETE CASCADE,
    photo_url TEXT NOT NULL,
    captured_at TIMESTAMP NOT NULL DEFAULT now()
);
CREATE TABLE store_photos (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    store_id UUID NOT NULL
        REFERENCES stores(id)
        ON DELETE CASCADE,
    photo_url TEXT NOT NULL,
    type TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);
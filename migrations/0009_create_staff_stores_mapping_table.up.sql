CREATE TABLE staff_stores_mapping (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    staff_id UUID NOT NULL
        REFERENCES staff_users(id)
        ON DELETE CASCADE,
    store_id UUID NOT NULL
        REFERENCES stores(id)
        ON DELETE CASCADE
);
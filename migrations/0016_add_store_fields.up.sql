ALTER TABLE stores
    ADD COLUMN owner_phone      VARCHAR(20)  NOT NULL DEFAULT '',
    ADD COLUMN additional_phone  VARCHAR(20),
    ADD COLUMN interior_photo    TEXT,
    ADD COLUMN exterior_photo    TEXT,
    ADD COLUMN comment           TEXT;

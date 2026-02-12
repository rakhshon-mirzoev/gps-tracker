ALTER TABLE stores
    DROP COLUMN IF EXISTS owner_phone,
    DROP COLUMN IF EXISTS additional_phone,
    DROP COLUMN IF EXISTS interior_photo,
    DROP COLUMN IF EXISTS exterior_photo,
    DROP COLUMN IF EXISTS comment;

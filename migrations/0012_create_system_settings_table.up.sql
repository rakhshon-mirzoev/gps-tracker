CREATE TABLE system_settings (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    key varchar(20) UNIQUE NOT NULL,
    value varchar(100) NOT NULL,
    description TEXT
);
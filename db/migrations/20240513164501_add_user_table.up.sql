CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE users (
    id UUID PRIMARY KEY NOT NULL DEFAULT uuid_generate_v4(),
    nip VARCHAR(20) NOT NULL,
    name VARCHAR(50) NOT NULL,
    password_hash VARCHAR(100),
    role VARCHAR(20) NOT NULL, -- IT, nurse, etc.
    identity_card_scan_img TEXT,
    access BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_users_nip ON users(nip);
CREATE INDEX idx_users_nip_role_createdat ON users(nip, role, created_at);
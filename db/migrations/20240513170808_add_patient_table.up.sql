CREATE TABLE patients (
    id UUID PRIMARY KEY NOT NULL DEFAULT uuid_generate_v4(),
    identity_number BIGINT NOT NULL UNIQUE,
    phone_number VARCHAR(15) NOT NULL,
    name VARCHAR(30) NOT NULL,
    birth_date DATE NOT NULL,
    gender VARCHAR(10) NOT NULL,
    identity_card_scan_img TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_patients_identity_number ON patients(identity_number);
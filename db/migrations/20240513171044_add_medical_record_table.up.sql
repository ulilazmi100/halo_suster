CREATE TABLE medical_records (
    id UUID PRIMARY KEY NOT NULL DEFAULT uuid_generate_v4(),
    patient_id UUID NOT NULL REFERENCES patients(id),
    symptoms TEXT NOT NULL,
    medications TEXT NOT NULL,
    created_by_nip BIGINT NOT NULL,
    created_by_name VARCHAR(50) NOT NULL,
    created_by_user_id UUID NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_medical_records_patient_id ON medical_records(patient_id);
CREATE INDEX idx_medical_records_created_by ON medical_records(created_by_nip, created_by_user_id);
CREATE INDEX idx_covering_medical_records ON medical_records(patient_id, symptoms, medications, created_by_nip, created_by_user_id);

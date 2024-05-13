DROP TABLE IF EXISTS medical_records CASCADE;

DROP INDEX IF EXISTS idx_patients_identity_number CASCADE;
DROP INDEX IF EXISTS idx_medical_records_created_by CASCADE;
DROP INDEX IF EXISTS idx_active_medical_records CASCADE;
DROP INDEX IF EXISTS idx_covering_medical_records CASCADE;
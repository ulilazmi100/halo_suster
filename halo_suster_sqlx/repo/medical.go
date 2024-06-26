package repo

import (
	"halo_suster/db/entities"
	"fmt"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
)

type MedicalRepo interface {
	GetPatient(patientIdentityNumber int64) (string, error)
	CreatePatient(patient *entities.PatientRegistrationPayload) error
	GetPatients(filter entities.GetPatientQueries) ([]entities.GetPatientResponse, error)
	CreateRecord(patient *entities.RecordRegistrationPayload, createdBy *entities.CreatedByDetail) error
	GetRecord(filter entities.GetRecordQueries) ([]entities.GetRecordResponse, error)
}

type medicalRepo struct {
	db *sqlx.DB
}

func NewMedicalRepo(db *sqlx.DB) MedicalRepo {
	return &medicalRepo{db}
}

func (r *medicalRepo) GetPatient(patientIdentityNumber int64) (string, error) {
	var identityNumber string
	query := "SELECT identity_number FROM patients WHERE identity_number = $1"

	row := r.db.QueryRow(query, patientIdentityNumber)
	err := row.Scan(&identityNumber)
	if err != nil {
		return "", err
	}

	return identityNumber, nil
}

func (r *medicalRepo) CreatePatient(patient *entities.PatientRegistrationPayload) error {
	statement := "INSERT INTO patients (identity_number, phone_number, name, birth_date, gender, identity_card_scan_img) VALUES ($1, $2, $3, $4, $5, $6)"

	_, err := r.db.Exec(statement, patient.IdentityNumber, patient.PhoneNumber, patient.Name, patient.BirthDate, patient.Gender, patient.IdentityCardScanImg)
	if err != nil {
		return err
	}

	return nil
}

func (r *medicalRepo) GetPatients(filter entities.GetPatientQueries) ([]entities.GetPatientResponse, error) {
	var patients []entities.GetPatientResponse
	var createdAt time.Time
	var birthDate time.Time
	query := "SELECT identity_number, phone_number, name, birth_date, gender, created_at FROM patients"

	query += getPatientConstructWhereQuery(filter)

	if filter.CreatedAt != "" {
		if filter.CreatedAt == "asc" {
			query += " ORDER BY created_at ASC"
		} else if filter.CreatedAt == "desc" {
			query += " ORDER BY created_at DESC"
		}
	} else {
		query += " ORDER BY created_at DESC"
	}

	query += " limit $1 offset $2"

	rows, err := r.db.Query(query, filter.Limit, filter.Offset)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		patient := entities.GetPatientResponse{}
		err := rows.Scan(&patient.IdentityNumber, &patient.PhoneNumber, &patient.Name, &birthDate, &patient.Gender, &createdAt)
		if err != nil {
			return nil, err
		}
		patient.BirthDate = birthDate.Format(time.RFC3339Nano)
		patient.CreatedAt = createdAt.Format(time.RFC3339Nano)
		patients = append(patients, patient)
	}

	return patients, nil
}

func (r *medicalRepo) CreateRecord(record *entities.RecordRegistrationPayload, createdBy *entities.CreatedByDetail) error {
	statement := "INSERT INTO medical_records (identity_number, symptoms, medications, created_by_nip, created_by_name, created_by_user_id) VALUES ($1, $2, $3, $4, $5, $6)"

	_, err := r.db.Exec(statement, record.IdentityNumber, record.Symptoms, record.Medications, entities.Int64ToString(createdBy.Nip), createdBy.Name, createdBy.UserId)
	if err != nil {
		return err
	}

	return nil
}

func (r *medicalRepo) GetRecord(filter entities.GetRecordQueries) ([]entities.GetRecordResponse, error) {
	var records []entities.GetRecordResponse
	var createdAt time.Time

	query := "SELECT identity_number, symptoms, medications, created_by_nip, created_by_name, created_by_user_id, created_at FROM medical_records"

	query += getRecordConstructWhereQuery(filter)

	if filter.CreatedAt != "" {
		if filter.CreatedAt == "asc" {
			query += " ORDER BY created_at ASC"
		} else if filter.CreatedAt == "desc" {
			query += " ORDER BY created_at DESC"
		}
	} else {
		query += " ORDER BY created_at DESC"
	}

	query += " limit $1 offset $2"

	rows, err := r.db.Query(query, filter.Limit, filter.Offset)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		record := entities.GetRecordResponse{}
		var identityNumber int64
		var birthDate time.Time
		var nipString string

		err := rows.Scan(&identityNumber, &record.Symptoms, &record.Medications, &nipString, &record.CreatedBy.Name, &record.CreatedBy.UserId, &createdAt)
		if err != nil {
			return nil, err
		}

		queryPatient := "SELECT phone_number, name, birth_date, gender, identity_card_scan_img FROM patients WHERE identity_number = $1"

		row := r.db.QueryRow(queryPatient, identityNumber)
		err = row.Scan(&record.IdentityDetail.PhoneNumber, &record.IdentityDetail.Name, &birthDate, &record.IdentityDetail.Gender, &record.IdentityDetail.IdentityCardScanImg)
		if err != nil {
			return []entities.GetRecordResponse{}, err
		}

		record.CreatedBy.Nip, err = entities.StringToInt64(nipString)
		if err != nil {
			return []entities.GetRecordResponse{}, err
		}
		record.IdentityDetail.IdentityNumber = identityNumber
		record.IdentityDetail.BirthDate = birthDate.Format(time.RFC3339Nano)

		record.CreatedAt = createdAt.Format(time.RFC3339Nano)

		records = append(records, record)
	}

	return records, nil
}

func getPatientConstructWhereQuery(filter entities.GetPatientQueries) string {
	whereSQL := []string{}

	if filter.IdentityNumber != nil {
		whereSQL = append(whereSQL, fmt.Sprintf(" identity_number = %d", *filter.IdentityNumber))
	}

	if filter.PhoneNumber != "" {
		whereSQL = append(whereSQL, " phone_number ILIKE '+"+filter.PhoneNumber+"%'")
	}

	if filter.Name != "" {
		whereSQL = append(whereSQL, " name ILIKE '%"+filter.Name+"%'")
	}

	if len(whereSQL) > 0 {
		return " WHERE " + strings.Join(whereSQL, " AND ")
	}

	return ""
}

func getRecordConstructWhereQuery(filter entities.GetRecordQueries) string {
	whereSQL := []string{}

	if filter.IdentityNumber != nil {
		whereSQL = append(whereSQL, fmt.Sprintf(" identity_number = %d", *filter.IdentityNumber))
	}

	if filter.CreatedByNip != nil {
		whereSQL = append(whereSQL, " created_by_nip = '"+entities.Int64ToString(*filter.CreatedByNip)+"'")
	}

	if filter.CreatedByUserId != "" {
		whereSQL = append(whereSQL, " created_by_user_id = '"+filter.CreatedByUserId+"'")
	}

	if len(whereSQL) > 0 {
		return " WHERE " + strings.Join(whereSQL, " AND ")
	}

	return ""
}

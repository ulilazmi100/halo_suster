package repo

import (
	"context"
	"fmt"
	"halo_suster/db/entities"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type MedicalRepo interface {
	GetPatient(ctx context.Context, patientIdentityNumber int64) (string, error)
	CreatePatient(ctx context.Context, patient *entities.PatientRegistrationPayload) error
	GetPatients(ctx context.Context, filter entities.GetPatientQueries) ([]entities.GetPatientResponse, error)
}

type medicalRepo struct {
	db *pgxpool.Pool
}

func NewMedicalRepo(db *pgxpool.Pool) MedicalRepo {
	return &medicalRepo{db}
}

func (r *medicalRepo) GetPatient(ctx context.Context, patientIdentityNumber int64) (string, error) {
	var identityNumber string
	query := "SELECT identity_number FROM patients WHERE identity_number = $1"

	row := r.db.QueryRow(ctx, query, patientIdentityNumber)
	err := row.Scan(&identityNumber)
	if err != nil {
		return "", err
	}

	return identityNumber, nil
}

func (r *medicalRepo) CreatePatient(ctx context.Context, patient *entities.PatientRegistrationPayload) error {
	statement := "INSERT INTO patients (identity_number, phone_number, name, birth_date, gender, identity_card_scan_img) VALUES ($1, $2, $3, $4, $5, $6)"

	_, err := r.db.Exec(ctx, statement, patient.IdentityNumber, patient.PhoneNumber, patient.Name, patient.BirthDate, patient.Gender, patient.IdentityCardScanImg)
	if err != nil {
		return err
	}

	return nil
}

func (r *medicalRepo) GetPatients(ctx context.Context, filter entities.GetPatientQueries) ([]entities.GetPatientResponse, error) {
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

	rows, err := r.db.Query(ctx, query, filter.Limit, filter.Offset)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		patient := entities.GetPatientResponse{}
		err := rows.Scan(&patient.IdentityNumber, &patient.PhoneNumber, &patient.Name, &birthDate, &patient.Gender, &createdAt)
		if err != nil {
			return nil, err
		}
		patient.BirthDate = birthDate.Format(time.RFC3339)
		patient.CreatedAt = createdAt.Format(time.RFC3339)
		patients = append(patients, patient)
	}

	return patients, nil
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

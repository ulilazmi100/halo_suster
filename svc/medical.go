package svc

import (
	"context"
	"halo_suster/db/entities"
	"halo_suster/repo"
	"halo_suster/responses"

	"github.com/jackc/pgx/v5"
)

type MedicalSvc interface {
	RegisterPatient(ctx context.Context, newPatient entities.PatientRegistrationPayload) error
	GetPatient(ctx context.Context, GetPatientQueries entities.GetPatientQueries) ([]entities.GetPatientResponse, error)
	RegisterRecord(ctx context.Context, newRecord entities.RecordRegistrationPayload, createdByDetail entities.CreatedByDetail) error
	GetRecord(ctx context.Context, GetRecordQueries entities.GetRecordQueries) ([]entities.GetRecordResponse, error)
}

type medicalSvc struct {
	repo repo.MedicalRepo
}

func NewMedicalSvc(repo repo.MedicalRepo) MedicalSvc {
	return &medicalSvc{repo}
}

func (s *medicalSvc) RegisterPatient(ctx context.Context, newPatient entities.PatientRegistrationPayload) error {
	if err := newPatient.Validate(); err != nil {
		return responses.NewBadRequestError(err.Error())
	}

	existingPatient, err := s.repo.GetPatient(ctx, newPatient.IdentityNumber)
	if err != nil {
		if err != pgx.ErrNoRows {
			return err
		}
	}

	if existingPatient != "" {
		return responses.NewConflictError("identityNumber exists")
	}

	err = s.repo.CreatePatient(ctx, &newPatient)
	if err != nil {
		return err
	}

	return nil
}

func (s *medicalSvc) GetPatient(ctx context.Context, GetPatientQueries entities.GetPatientQueries) ([]entities.GetPatientResponse, error) {

	if GetPatientQueries.IdentityNumber != nil && entities.ValidateIdentityNumber(*GetPatientQueries.IdentityNumber) != nil {
		return nil, nil
	}

	patients, err := s.repo.GetPatients(ctx, GetPatientQueries)
	if err != nil {
		if err == pgx.ErrNoRows {
			return []entities.GetPatientResponse{}, nil
		}
		return nil, err
	}

	return patients, nil
}

func (s *medicalSvc) RegisterRecord(ctx context.Context, newRecord entities.RecordRegistrationPayload, createdByDetail entities.CreatedByDetail) error {
	if err := newRecord.Validate(); err != nil {
		return responses.NewBadRequestError(err.Error())
	}

	existingPatient, err := s.repo.GetPatient(ctx, newRecord.IdentityNumber)
	if err != nil {
		if err == pgx.ErrNoRows {
			return responses.NewBadRequestError("identityNumber is not exist")
		}
	}

	if existingPatient == "" {
		return responses.NewBadRequestError("identityNumber is not exist")
	}

	err = s.repo.CreateRecord(ctx, &newRecord, &createdByDetail)
	if err != nil {
		return err
	}

	return nil
}

func (s *medicalSvc) GetRecord(ctx context.Context, GetRecordQueries entities.GetRecordQueries) ([]entities.GetRecordResponse, error) {

	if GetRecordQueries.IdentityNumber != nil && entities.ValidateIdentityNumber(*GetRecordQueries.IdentityNumber) != nil {
		return nil, nil
	}

	patients, err := s.repo.GetRecord(ctx, GetRecordQueries)
	if err != nil {
		if err == pgx.ErrNoRows {
			return []entities.GetRecordResponse{}, nil
		}
		return nil, err
	}

	return patients, nil
}

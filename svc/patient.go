package svc

import (
	"context"
	"halo_suster/db/entities"
	"halo_suster/repo"
	"halo_suster/responses"

	"github.com/jackc/pgx/v5"
)

type PatientSvc interface {
	RegisterPatient(ctx context.Context, newPatient entities.PatientRegistrationPayload) error
	GetPatient(ctx context.Context, GetPatientQueries entities.GetPatientQueries) ([]entities.GetPatientResponse, error)
}

type patientSvc struct {
	repo repo.PatientRepo
}

func NewPatientSvc(repo repo.PatientRepo) PatientSvc {
	return &patientSvc{repo}
}

func (s *patientSvc) RegisterPatient(ctx context.Context, newPatient entities.PatientRegistrationPayload) error {
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

func (s *patientSvc) GetPatient(ctx context.Context, GetPatientQueries entities.GetPatientQueries) ([]entities.GetPatientResponse, error) {

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

package svc

import (
	"halo_suster/db/entities"
	"halo_suster/repo"
	"halo_suster/responses"
	"database/sql"
)

type MedicalSvc interface {
	RegisterPatient(newPatient entities.PatientRegistrationPayload) error
	GetPatient(GetPatientQueries entities.GetPatientQueries) ([]entities.GetPatientResponse, error)
	RegisterRecord(newRecord entities.RecordRegistrationPayload, createdByDetail entities.CreatedByDetail) error
	GetRecord(GetRecordQueries entities.GetRecordQueries) ([]entities.GetRecordResponse, error)
}

type medicalSvc struct {
	repo repo.MedicalRepo
}

func NewMedicalSvc(repo repo.MedicalRepo) MedicalSvc {
	return &medicalSvc{repo}
}

func (s *medicalSvc) RegisterPatient(newPatient entities.PatientRegistrationPayload) error {
	if err := newPatient.Validate(); err != nil {
		return responses.NewBadRequestError(err.Error())
	}

	existingPatient, err := s.repo.GetPatient(newPatient.IdentityNumber)
	if err != nil {
		if err != sql.ErrNoRows {
			return err
		}
	}

	if existingPatient != "" {
		return responses.NewConflictError("identityNumber exists")
	}

	err = s.repo.CreatePatient(&newPatient)
	if err != nil {
		return err
	}

	return nil
}

func (s *medicalSvc) GetPatient(GetPatientQueries entities.GetPatientQueries) ([]entities.GetPatientResponse, error) {

	if GetPatientQueries.IdentityNumber != nil && entities.ValidateIdentityNumber(*GetPatientQueries.IdentityNumber) != nil {
		return nil, nil
	}

	patients, err := s.repo.GetPatients(GetPatientQueries)
	if err != nil {
		if err == sql.ErrNoRows {
			return []entities.GetPatientResponse{}, nil
		}
		return nil, err
	}

	return patients, nil
}

func (s *medicalSvc) RegisterRecord(newRecord entities.RecordRegistrationPayload, createdByDetail entities.CreatedByDetail) error {
	if err := newRecord.Validate(); err != nil {
		return responses.NewBadRequestError(err.Error())
	}

	existingPatient, err := s.repo.GetPatient(newRecord.IdentityNumber)
	if err != nil {
		if err == sql.ErrNoRows {
			return responses.NewNotFoundError("identityNumber is not exist")
		}
	}

	if existingPatient == "" {
		return responses.NewNotFoundError("identityNumber is not exist")
	}

	err = s.repo.CreateRecord(&newRecord, &createdByDetail)
	if err != nil {
		return err
	}

	return nil
}

func (s *medicalSvc) GetRecord(GetRecordQueries entities.GetRecordQueries) ([]entities.GetRecordResponse, error) {

	if GetRecordQueries.IdentityNumber != nil && entities.ValidateIdentityNumber(*GetRecordQueries.IdentityNumber) != nil {
		return nil, nil
	}

	patients, err := s.repo.GetRecord(GetRecordQueries)
	if err != nil {
		if err == sql.ErrNoRows {
			return []entities.GetRecordResponse{}, nil
		}
		return nil, err
	}

	return patients, nil
}

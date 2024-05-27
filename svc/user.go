package svc

import (
	"halo_suster/crypto"
	"halo_suster/db/entities"
	"halo_suster/repo"
	"halo_suster/responses"
	"database/sql"
	"strings"
)

type UserSvc interface {
	Register(newUser entities.RegistrationPayload) (string, string, error)
	Login(user entities.Credential) (string, string, string, error)
	NurseRegister(newUser entities.NurseRegistrationPayload) (string, error)
	NurseLogin(creds entities.Credential) (string, string, string, error)
	UpdateNurse(nurseId string, updatePayload entities.NurseUpdatePayload) error
	DeleteNurse(nurseId string) error
	AccessNurse(nurseId string, password entities.NurseAccessPayload) error
	GetUser(GetUserQueries entities.GetUserQueries) ([]entities.GetUserResponse, error)
}

type userSvc struct {
	repo repo.UserRepo
}

func NewUserSvc(repo repo.UserRepo) UserSvc {
	return &userSvc{repo}
}

func (s *userSvc) Register(newUser entities.RegistrationPayload) (string, string, error) {
	if err := newUser.Validate(); err != nil {
		return "", "", responses.NewBadRequestError(err.Error())
	}

	existingUser, err := s.repo.GetUser(entities.Int64ToString(newUser.Nip))
	if err != nil {
		if err != sql.ErrNoRows {
			return "", "", err
		}
	}

	if existingUser != nil {
		return "", "", responses.NewConflictError("user already exists")
	}

	hashedPassword, err := crypto.GenerateHashedPassword(newUser.Password)
	if err != nil {
		return "", "", err
	}

	id, err := s.repo.CreateUser(&newUser, hashedPassword)
	if err != nil {
		return "", "", err
	}

	token, err := crypto.GenerateToken(id, entities.Int64ToString(newUser.Nip), newUser.Name)
	if err != nil {
		return "", "", err
	}

	return id, token, nil
}

func (s *userSvc) Login(creds entities.Credential) (string, string, string, error) {
	if strings.HasPrefix(creds.Nip, "303") {
		return "", "", "", responses.NewNotFoundError("user is not from IT (nip not starts with 615)")
	}

	err := entities.ValidateITStaffNipFormat(creds.Nip)

	if err != nil {
		return "", "", "", responses.NewBadRequestError(err.Error())
	}

	if err := creds.Validate(); err != nil {
		return "", "", "", responses.NewBadRequestError(err.Error())
	}

	user, err := s.repo.GetUser(creds.Nip)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", "", "", responses.NewNotFoundError("user not found")
		}
		return "", "", "", err
	}

	err = crypto.VerifyPassword(creds.Password, user.Password)
	if err != nil {
		return "", "", "", responses.NewBadRequestError("wrong password!")
	}

	token, err := crypto.GenerateToken(user.Id, user.Nip, user.Name)
	if err != nil {
		return "", "", "", responses.NewBadRequestError(err.Error())
	}

	return user.Id, user.Name, token, nil
}

func (s *userSvc) NurseRegister(newUser entities.NurseRegistrationPayload) (string, error) {
	if err := newUser.Validate(); err != nil {
		return "", responses.NewBadRequestError(err.Error())
	}

	existingUser, err := s.repo.GetUser(entities.Int64ToString(newUser.Nip))
	if err != nil {
		if err != sql.ErrNoRows {
			return "", err
		}
	}

	if existingUser != nil {
		return "", responses.NewConflictError("user already exists")
	}

	id, err := s.repo.CreateNurseUser(&newUser)
	if err != nil {
		return "", err
	}

	return id, nil
}

func (s *userSvc) NurseLogin(creds entities.Credential) (string, string, string, error) {
	if strings.HasPrefix(creds.Nip, "615") {
		return "", "", "", responses.NewNotFoundError("user is not from nurse (nip not starts with 303)")
	}

	err := entities.ValidateNurseNipFormat(creds.Nip)

	if err != nil {
		return "", "", "", responses.NewBadRequestError(err.Error())
	}

	if err := creds.Validate(); err != nil {
		return "", "", "", responses.NewBadRequestError(err.Error())
	}

	user, err := s.repo.GetUser(creds.Nip)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", "", "", responses.NewNotFoundError("user not found")
		}
		return "", "", "", err
	}

	if !user.Access {
		return "", "", "", responses.NewBadRequestError("no access!")
	}

	err = crypto.VerifyPassword(creds.Password, user.Password)
	if err != nil {
		return "", "", "", responses.NewBadRequestError("wrong password!")
	}

	token, err := crypto.GenerateToken(user.Id, user.Nip, user.Name)
	if err != nil {
		return "", "", "", responses.NewBadRequestError(err.Error())
	}

	return user.Id, user.Name, token, nil
}

func (s *userSvc) UpdateNurse(nurseId string, updatePayload entities.NurseUpdatePayload) error {

	err := updatePayload.Validate()

	if err != nil {
		return responses.NewBadRequestError(err.Error())
	}

	err = entities.ValidateNurseNipFormat(updatePayload.Nip)

	if err != nil {
		return responses.NewBadRequestError(err.Error())
	}

	if !entities.IsValidUUID(nurseId) {
		return responses.NewNotFoundError("user not found")
	}

	existingUser, err := s.repo.GetUser(entities.Int64ToString(updatePayload.Nip))
	if err != nil {
		if err != sql.ErrNoRows {
			return err
		}
	}

	if existingUser != nil {
		return responses.NewConflictError("conflict nip already used")
	}

	user, err := s.repo.GetUserNipById(nurseId)
	if err != nil {
		if err == sql.ErrNoRows {
			return responses.NewNotFoundError("user not found")
		}
		return err
	}

	if !strings.HasPrefix(user.Nip, "303") {
		return responses.NewNotFoundError("user is not a nurse (nip not starts with 303)")
	}

	res, err := s.repo.UpdateNurse(nurseId, updatePayload)

	rowAffected, err := res.RowsAffected()
	if rowAffected == 0 {
		return responses.NewNotFoundError(err.Error())
	}

	if err != nil {
		return responses.NewInternalServerError(err.Error())
	}

	return err
}

func (s *userSvc) DeleteNurse(nurseId string) error {
	if !entities.IsValidUUID(nurseId) {
		return responses.NewNotFoundError("user not found")
	}

	user, err := s.repo.GetUserNipById(nurseId)
	if err != nil {
		if err == sql.ErrNoRows {
			return responses.NewNotFoundError("user not found")
		}
		return err
	}

	if !strings.HasPrefix(user.Nip, "303") {
		return responses.NewNotFoundError("user is not a nurse (nip not starts with 303)")
	}

	res, err := s.repo.DeleteNurse(nurseId)

	rowAffected, err := res.RowsAffected()
	if rowAffected == 0 {
		return responses.NewNotFoundError(err.Error())
	}

	if err != nil {
		return responses.NewInternalServerError(err.Error())
	}

	return err
}

func (s *userSvc) AccessNurse(nurseId string, accessPayload entities.NurseAccessPayload) error {
	err := accessPayload.Validate()

	if err != nil {
		return responses.NewBadRequestError(err.Error())
	}

	if !entities.IsValidUUID(nurseId) {
		return responses.NewNotFoundError("user not found")
	}

	user, err := s.repo.GetUserNipById(nurseId)
	if err != nil {
		if err == sql.ErrNoRows {
			return responses.NewNotFoundError("user not found")
		}
		return err
	}

	if !strings.HasPrefix(user.Nip, "303") {
		return responses.NewNotFoundError("user is not a nurse (nip not starts with 303)")
	}

	hashedPassword, err := crypto.GenerateHashedPassword(accessPayload.Password)
	if err != nil {
		return err
	}

	res, err := s.repo.UpdateAccessNurse(nurseId, hashedPassword)

	rowAffected, err := res.RowsAffected()
	if rowAffected == 0 {
		return responses.NewNotFoundError(err.Error())
	}

	if err != nil {
		return responses.NewInternalServerError(err.Error())
	}

	return err
}

func (s *userSvc) GetUser(GetUserQueries entities.GetUserQueries) ([]entities.GetUserResponse, error) {

	if GetUserQueries.UserId != "" && !entities.IsValidUUID(GetUserQueries.UserId) {
		return nil, nil
	}

	users, err := s.repo.GetUsers(GetUserQueries)
	if err != nil {
		if err == sql.ErrNoRows {
			return []entities.GetUserResponse{}, nil
		}
		return nil, err
	}

	return users, nil
}

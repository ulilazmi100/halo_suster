package svc

import (
	"context"
	"halo_suster/crypto"
	"halo_suster/db/entities"
	"halo_suster/repo"
	"halo_suster/responses"
	"strings"

	"github.com/jackc/pgx/v5"
)

type UserSvc interface {
	Register(ctx context.Context, newUser entities.RegistrationPayload) (string, string, error)
	Login(ctx context.Context, user entities.Credential) (string, string, string, error)
	NurseRegister(ctx context.Context, newUser entities.NurseRegistrationPayload) (string, error)
	NurseLogin(ctx context.Context, creds entities.Credential) (string, string, string, error)
	UpdateNurse(ctx context.Context, nurseId string, updatePayload entities.NurseUpdatePayload) error
	DeleteNurse(ctx context.Context, nurseId string) error
	AccessNurse(ctx context.Context, nurseId string, password entities.NurseAccessPayload) error
	GetUser(ctx context.Context, GetUserQueries entities.GetUserQueries) ([]entities.GetUserResponse, error)
}

type userSvc struct {
	repo repo.UserRepo
}

func NewUserSvc(repo repo.UserRepo) UserSvc {
	return &userSvc{repo}
}

func (s *userSvc) Register(ctx context.Context, newUser entities.RegistrationPayload) (string, string, error) {
	if err := newUser.Validate(); err != nil {
		return "", "", responses.NewBadRequestError(err.Error())
	}

	existingUser, err := s.repo.GetUser(ctx, newUser.Nip)
	if err != nil {
		if err != pgx.ErrNoRows {
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

	// // Use a transaction for creating the user and generating the token
	// tx, err := s.repo.BeginTx(ctx)
	// if err != nil {
	// 	return "", "", err
	// }
	// defer tx.Rollback(ctx)

	id, err := s.repo.CreateUser(ctx, &newUser, hashedPassword)
	if err != nil {
		return "", "", err
	}

	token, err := crypto.GenerateToken(id, newUser.Nip, newUser.Name)
	if err != nil {
		return "", "", err
	}

	// if err := tx.Commit(ctx); err != nil {
	// 	return "", "", err
	// }

	return id, token, nil
}

func (s *userSvc) Login(ctx context.Context, creds entities.Credential) (string, string, string, error) {
	err := entities.ValidateITStaffNipFormat(creds.Nip)

	if err != nil {
		return "", "", "", responses.NewNotFoundError("User not found")
	}

	if err := creds.Validate(); err != nil {
		return "", "", "", responses.NewBadRequestError(err.Error())
	}

	user, err := s.repo.GetUser(ctx, creds.Nip)
	if err != nil {
		if err == pgx.ErrNoRows {
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

func (s *userSvc) NurseRegister(ctx context.Context, newUser entities.NurseRegistrationPayload) (string, error) {
	if err := newUser.Validate(); err != nil {
		return "", responses.NewBadRequestError(err.Error())
	}

	existingUser, err := s.repo.GetUser(ctx, newUser.Nip)
	if err != nil {
		if err != pgx.ErrNoRows {
			return "", err
		}
	}

	if existingUser != nil {
		return "", responses.NewConflictError("user already exists")
	}

	id, err := s.repo.CreateNurseUser(ctx, &newUser)
	if err != nil {
		return "", err
	}

	return id, nil
}

func (s *userSvc) NurseLogin(ctx context.Context, creds entities.Credential) (string, string, string, error) {
	err := entities.ValidateNurseNipFormat(creds.Nip)

	if err != nil {
		return "", "", "", responses.NewNotFoundError("User not found")
	}

	if err := creds.Validate(); err != nil {
		return "", "", "", responses.NewBadRequestError(err.Error())
	}

	user, err := s.repo.GetUser(ctx, creds.Nip)
	if err != nil {
		if err == pgx.ErrNoRows {
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

func (s *userSvc) UpdateNurse(ctx context.Context, nurseId string, updatePayload entities.NurseUpdatePayload) error {
	err := updatePayload.Validate()

	if err != nil {
		return responses.NewBadRequestError(err.Error())
	}

	if !entities.IsValidUUID(nurseId) {
		return responses.NewNotFoundError("user not found")
	}

	existingUser, err := s.repo.GetUser(ctx, updatePayload.Nip)
	if err != nil {
		if err != pgx.ErrNoRows {
			return err
		}
	}

	if existingUser != nil {
		return responses.NewConflictError("conflict nip already used")
	}

	user, err := s.repo.GetUserNipById(ctx, nurseId)
	if err != nil {
		if err == pgx.ErrNoRows {
			return responses.NewNotFoundError("user not found")
		}
		return err
	}

	if !strings.HasPrefix(user.Nip, "303") {
		return responses.NewNotFoundError("user is not a nurse (nip not starts with 303)")
	}

	res, err := s.repo.UpdateNurse(ctx, nurseId, updatePayload)

	if res.RowsAffected() == 0 {
		return responses.NewNotFoundError(err.Error())
	}

	if err != nil {
		return responses.NewInternalServerError(err.Error())
	}

	return err
}

func (s *userSvc) DeleteNurse(ctx context.Context, nurseId string) error {
	if !entities.IsValidUUID(nurseId) {
		return responses.NewNotFoundError("user not found")
	}

	user, err := s.repo.GetUserNipById(ctx, nurseId)
	if err != nil {
		if err == pgx.ErrNoRows {
			return responses.NewNotFoundError("user not found")
		}
		return err
	}

	if !strings.HasPrefix(user.Nip, "303") {
		return responses.NewNotFoundError("user is not a nurse (nip not starts with 303)")
	}

	res, err := s.repo.DeleteNurse(ctx, nurseId)

	if res.RowsAffected() == 0 {
		return responses.NewNotFoundError(err.Error())
	}

	if err != nil {
		return responses.NewInternalServerError(err.Error())
	}

	return err
}

func (s *userSvc) AccessNurse(ctx context.Context, nurseId string, accessPayload entities.NurseAccessPayload) error {
	err := accessPayload.Validate()

	if err != nil {
		return responses.NewBadRequestError(err.Error())
	}

	if !entities.IsValidUUID(nurseId) {
		return responses.NewNotFoundError("user not found")
	}

	user, err := s.repo.GetUserNipById(ctx, nurseId)
	if err != nil {
		if err == pgx.ErrNoRows {
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

	res, err := s.repo.UpdateAccessNurse(ctx, nurseId, hashedPassword)

	if res.RowsAffected() == 0 {
		return responses.NewNotFoundError(err.Error())
	}

	if err != nil {
		return responses.NewInternalServerError(err.Error())
	}

	return err
}

func (s *userSvc) GetUser(ctx context.Context, GetUserQueries entities.GetUserQueries) ([]entities.GetUserResponse, error) {

	if GetUserQueries.UserId != "" && !entities.IsValidUUID(GetUserQueries.UserId) {
		return nil, nil
	}

	users, err := s.repo.GetUsers(ctx, GetUserQueries)
	if err != nil {
		if err == pgx.ErrNoRows {
			return []entities.GetUserResponse{}, nil
		}
		return nil, err
	}

	return users, nil
}

package repo

import (
	"context"
	"halo_suster/db/entities"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepo interface {
	GetUser(ctx context.Context, nip string) (*entities.User, error)
	CreateUser(ctx context.Context, user *entities.RegistrationPayload, hashPassword string) (string, error)
	CreateUserTx(ctx context.Context, tx pgx.Tx, user *entities.RegistrationPayload, hashPassword string) (string, error)
	CreateNurseUser(ctx context.Context, user *entities.NurseRegistrationPayload) (string, error)
	GetUserNipById(ctx context.Context, id string) (*entities.User, error)
	UpdateNurse(ctx context.Context, nurseId string, updatePayload entities.NurseUpdatePayload) (pgconn.CommandTag, error)
	UpdateAccessNurse(ctx context.Context, nurseId string, passwordHash string) (pgconn.CommandTag, error)
	DeleteNurse(ctx context.Context, userId string) (pgconn.CommandTag, error)
	GetUsers(ctx context.Context, filter entities.GetUserQueries) ([]entities.GetUserResponse, error)
	BeginTx(ctx context.Context) (pgx.Tx, error)
}

type userRepo struct {
	db *pgxpool.Pool
}

func NewUserRepo(db *pgxpool.Pool) UserRepo {
	return &userRepo{db}
}

func (r *userRepo) GetUser(ctx context.Context, nip string) (*entities.User, error) {
	var user entities.User
	query := "SELECT id, name, password_hash, access FROM users WHERE nip = $1"

	row := r.db.QueryRow(ctx, query, nip)
	err := row.Scan(&user.Id, &user.Name, &user.Password, &user.Access)
	if err != nil {
		return nil, err
	}
	user.Nip = nip

	return &user, nil
}

func (r *userRepo) CreateUser(ctx context.Context, user *entities.RegistrationPayload, hashPassword string) (string, error) {
	var id string
	role := CheckRoleForRegister(user.Nip)
	statement := "INSERT INTO users (name, nip, password_hash, role) VALUES ($1, $2, $3, $4) RETURNING id"

	row := r.db.QueryRow(ctx, statement, user.Name, user.Nip, hashPassword, role)
	if err := row.Scan(&id); err != nil {
		return "", err
	}

	return id, nil
}

func (r *userRepo) CreateUserTx(ctx context.Context, tx pgx.Tx, user *entities.RegistrationPayload, hashPassword string) (string, error) {
	var id string
	statement := "INSERT INTO users (name, nip, password_hash) VALUES ($1, $2, $3) RETURNING id"

	row := tx.QueryRow(ctx, statement, user.Name, user.Nip, hashPassword)
	if err := row.Scan(&id); err != nil {
		return "", err
	}

	return id, nil
}

func (r *userRepo) CreateNurseUser(ctx context.Context, user *entities.NurseRegistrationPayload) (string, error) {
	var id string
	role := CheckRoleForRegister(user.Nip)
	statement := "INSERT INTO users (name, nip, password_hash, role, access) VALUES ($1, $2, '', $3, false) RETURNING id"

	row := r.db.QueryRow(ctx, statement, user.Name, user.Nip, role)
	if err := row.Scan(&id); err != nil {
		return "", err
	}

	return id, nil
}

func (r *userRepo) GetUserNipById(ctx context.Context, id string) (*entities.User, error) {
	var user entities.User
	query := "SELECT nip FROM users WHERE id = $1"

	row := r.db.QueryRow(ctx, query, id)
	err := row.Scan(&user.Nip)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *userRepo) UpdateNurse(ctx context.Context, nurseId string, updatePayload entities.NurseUpdatePayload) (pgconn.CommandTag, error) {
	statement := "UPDATE users SET nip = $1, name = $2 WHERE id = $3"

	res, err := r.db.Exec(ctx, statement, updatePayload.Nip, updatePayload.Name, nurseId)

	return res, err
}

func (r *userRepo) DeleteNurse(ctx context.Context, userId string) (pgconn.CommandTag, error) {
	statement := "DELETE FROM users WHERE id = $1"

	res, err := r.db.Exec(ctx, statement, userId)
	return res, err
}

func (r *userRepo) UpdateAccessNurse(ctx context.Context, nurseId string, passwordHash string) (pgconn.CommandTag, error) {
	statement := "UPDATE users SET password_hash = $1, access = true WHERE id = $2"

	res, err := r.db.Exec(ctx, statement, passwordHash, nurseId)

	return res, err
}

func (r *userRepo) GetUsers(ctx context.Context, filter entities.GetUserQueries) ([]entities.GetUserResponse, error) {
	var users []entities.GetUserResponse
	var createdAt time.Time
	query := "SELECT id, nip, name, created_at FROM users"

	query += getUserConstructWhereQuery(filter)

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
		user := entities.GetUserResponse{}
		err := rows.Scan(&user.UserId, &user.Nip, &user.Name, &createdAt)
		if err != nil {
			return nil, err
		}
		user.CreatedAt = createdAt.Format(time.RFC3339)
		users = append(users, user)
	}

	return users, nil
}

func (r *userRepo) BeginTx(ctx context.Context) (pgx.Tx, error) {
	return r.db.Begin(ctx)
}

func CheckRoleForRegister(nip string) string {
	if strings.HasPrefix(nip, "615") {
		return "it"
	} else if strings.HasPrefix(nip, "303") {
		return "nurse"
	}
	return ""
}

func getUserConstructWhereQuery(filter entities.GetUserQueries) string {
	whereSQL := []string{}

	if filter.UserId != "" {
		whereSQL = append(whereSQL, " id = '"+filter.UserId+"'")
	}

	if filter.Name != "" {
		whereSQL = append(whereSQL, " name ILIKE '%"+filter.Name+"%'")
	}

	if filter.Nip != "" {
		whereSQL = append(whereSQL, " nip ILIKE '"+filter.Nip+"%'")
	}

	if filter.Role == "it" {
		whereSQL = append(whereSQL, " role = '"+filter.Role+"'")
	}

	if filter.Role == "nurse" {
		whereSQL = append(whereSQL, " role = '"+filter.Role+"'")
	}

	if len(whereSQL) > 0 {
		return " WHERE " + strings.Join(whereSQL, " AND ")
	}

	return ""
}

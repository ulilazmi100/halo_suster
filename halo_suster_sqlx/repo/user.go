package repo

import (
	"halo_suster/db/entities"
	"database/sql"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
)

type UserRepo interface {
	GetUser(nip string) (*entities.User, error)
	CreateUser(user *entities.RegistrationPayload, hashPassword string) (string, error)
	CreateNurseUser(user *entities.NurseRegistrationPayload) (string, error)
	GetUserNipById(id string) (*entities.User, error)
	UpdateNurse(nurseId string, updatePayload entities.NurseUpdatePayload) (sql.Result, error)
	UpdateAccessNurse(nurseId string, passwordHash string) (sql.Result, error)
	DeleteNurse(userId string) (sql.Result, error)
	GetUsers(filter entities.GetUserQueries) ([]entities.GetUserResponse, error)
}

type userRepo struct {
	db *sqlx.DB
}

func NewUserRepo(db *sqlx.DB) UserRepo {
	return &userRepo{db}
}

func (r *userRepo) GetUser(nip string) (*entities.User, error) {
	var user entities.User
	query := "SELECT id, name, password_hash, access FROM users WHERE nip = $1"

	row := r.db.QueryRow(query, nip)
	err := row.Scan(&user.Id, &user.Name, &user.Password, &user.Access)
	if err != nil {
		return nil, err
	}
	user.Nip = nip

	return &user, nil
}

func (r *userRepo) CreateUser(user *entities.RegistrationPayload, hashPassword string) (string, error) {
	var id string
	role := CheckRoleForRegister(entities.Int64ToString(user.Nip))
	statement := "INSERT INTO users (name, nip, password_hash, role) VALUES ($1, $2, $3, $4) RETURNING id"

	row := r.db.QueryRow(statement, user.Name, entities.Int64ToString(user.Nip), hashPassword, role)
	if err := row.Scan(&id); err != nil {
		return "", err
	}

	return id, nil
}

func (r *userRepo) CreateNurseUser(user *entities.NurseRegistrationPayload) (string, error) {
	var id string
	role := CheckRoleForRegister(entities.Int64ToString(user.Nip))
	statement := "INSERT INTO users (name, nip, password_hash, role, access) VALUES ($1, $2, '', $3, false) RETURNING id"

	row := r.db.QueryRow(statement, user.Name, entities.Int64ToString(user.Nip), role)
	if err := row.Scan(&id); err != nil {
		return "", err
	}

	return id, nil
}

func (r *userRepo) GetUserNipById(id string) (*entities.User, error) {
	var user entities.User
	query := "SELECT nip FROM users WHERE id = $1"

	row := r.db.QueryRow(query, id)
	err := row.Scan(&user.Nip)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *userRepo) UpdateNurse(nurseId string, updatePayload entities.NurseUpdatePayload) (sql.Result, error) {
	statement := "UPDATE users SET nip = $1, name = $2 WHERE id = $3"

	res, err := r.db.Exec(statement, entities.Int64ToString(updatePayload.Nip), updatePayload.Name, nurseId)

	return res, err
}

func (r *userRepo) DeleteNurse(userId string) (sql.Result, error) {
	statement := "DELETE FROM users WHERE id = $1"

	res, err := r.db.Exec(statement, userId)
	return res, err
}

func (r *userRepo) UpdateAccessNurse(nurseId string, passwordHash string) (sql.Result, error) {
	statement := "UPDATE users SET password_hash = $1, access = true WHERE id = $2"

	res, err := r.db.Exec(statement, passwordHash, nurseId)

	return res, err
}

func (r *userRepo) GetUsers(filter entities.GetUserQueries) ([]entities.GetUserResponse, error) {
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

	rows, err := r.db.Query(query, filter.Limit, filter.Offset)
	if err != nil {
		return nil, err
	}

	var nip string
	for rows.Next() {
		user := entities.GetUserResponse{}
		err := rows.Scan(&user.UserId, &nip, &user.Name, &createdAt)
		if err != nil {
			return nil, err
		}
		user.Nip, err = entities.StringToInt64(nip)
		if err != nil {
			return nil, err
		}
		user.CreatedAt = createdAt.Format(time.RFC3339Nano)
		users = append(users, user)
	}

	return users, nil
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

	if filter.Nip != nil {
		whereSQL = append(whereSQL, " nip ILIKE '"+entities.Int64ToString(*filter.Nip)+"%'")
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

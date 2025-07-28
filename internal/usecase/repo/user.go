package repo

import (
	"context"
	"database/sql"
	"errors"
	"strconv"
	"strings"
	"time"

	"shop/config"
	"shop/internal/controller/http/token"
	"shop/internal/entity"
	"shop/pkg/logger"
	"shop/pkg/postgres"
)

type UserRepo struct {
	pg     *postgres.Postgres
	config *config.Config
	logger *logger.Logger
}

// New -.
func NewUserRepo(pg *postgres.Postgres, config *config.Config, logger *logger.Logger) *UserRepo {
	return &UserRepo{
		pg:     pg,
		config: config,
		logger: logger,
	}
}

func (r *UserRepo) Create(ctx context.Context, req *entity.CreateUser) error {
	query := `
		INSERT INTO users (
			name,
			password,
			phone_number
		) VALUES($1, $2, $3)`

	_, err := r.pg.Pool.Exec(ctx, query,
		req.Name,
		req.Password,
		req.PhoneNumber,
	)
	if err != nil {
		return err
	}

	return nil
}

func (r *UserRepo) Login(ctx context.Context, req *entity.LoginReq) (*entity.LoginRes, error) {
	query := `
		SELECT
			password,
			id,
			role
		FROM
			users
		WHERE
			phone_number = $1
		AND
			deleted_at = 0
	`
	row := r.pg.Pool.QueryRow(ctx, query, req.Login)
	var password string
	var id string
	var role string
	err := row.Scan(&password, &id, &role)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	if password != req.Password {
		return nil, errors.New("invalid password")
	}

	token := token.GenerateJWTToken(id, role)

	return &entity.LoginRes{Token: token.AccessToken,
		Message: "success"}, nil
}

func (r *UserRepo) GetById(ctx context.Context, req *entity.ById) (*entity.UserInfo, error) {

	var res entity.UserInfo
	var createdAt time.Time

	query := `
	SELECT
		id,
		name,
		phone_number,
		debt,
		created_at
	FROM 
		users
	WHERE 
		deleted_at = 0
	AND 
		id = $1
	`

	row := r.pg.Pool.QueryRow(ctx, query, req.Id)
	err := row.Scan(
		&res.ID,
		&res.Name,
		&res.PhoneNumber,
		&res.Debt,
		&createdAt,
	)
	if err != nil {
		return nil, err
	}
	res.CreatedAt = createdAt.Format("2006-01-02 15:04:05")

	return &res, nil
}

func (r *UserRepo) GetAll(ctx context.Context, req *entity.Filter) (*entity.UserList, error) {

	resp := &entity.UserList{}

	query := `
	SELECT
		COUNT(id) OVER () AS total_count,
		id,
		name,
		phone_number,
		debt,
		created_at
	FROM
		users
	WHERE
		deleted_at = 0
	`

	var args []interface{}

	if req.Limit == 0 {
		query += " OFFSET $1"
		args = append(args, req.Offset)
	} else {
		query += " LIMIT $1 OFFSET $2"
		args = append(args, req.Limit)
		args = append(args, req.Offset)
	}
	rows, err := r.pg.Pool.Query(ctx, query, args...)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user list is empty")
		}
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		res := entity.UserInfo{}
		var count int
		var createdAt time.Time

		err := rows.Scan(
			&count,
			&res.ID,
			&res.Name,
			&res.PhoneNumber,
			&res.Debt,
			&createdAt,
		)
		if err != nil {
			return nil, err
		}
		res.CreatedAt = createdAt.Format("2006-01-02 15:04:05")

		resp.Users = append(resp.Users, res)
		resp.Count = count
	}

	return resp, nil
}

func (r *UserRepo) Update(ctx context.Context, req *entity.UpdateUser) error {
	query := `
	UPDATE
		users
	SET`

	var conditions []string
	var args []interface{}

	if req.Name != "" && req.Name != "string" {
		conditions = append(conditions, " name = $"+strconv.Itoa(len(args)+1))
		args = append(args, req.Name)
	}
	if req.PhoneNumber != "" && req.PhoneNumber != "string" {
		conditions = append(conditions, " phone_number = $"+strconv.Itoa(len(args)+1))
		args = append(args, req.PhoneNumber)
	}

	conditions = append(conditions, " updated_at = now()")
	query += strings.Join(conditions, ", ")
	query += " WHERE id = $" + strconv.Itoa(len(args)+1) + " AND deleted_at = 0"

	args = append(args, req.Id)

	_, err := r.pg.Pool.Exec(ctx, query, args...)
	if err != nil {
		return err
	}

	return nil
}

func (r *UserRepo) Delete(ctx context.Context, req *entity.ById) error {

	_, err := r.pg.Pool.Exec(ctx, `UPDATE users SET deleted_at = EXTRACT(EPOCH FROM NOW()) WHERE id = $1`, req.Id)
	if err != nil {
		return err
	}

	return nil
}

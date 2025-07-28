package repo

import (
	"context"
	"database/sql"
	"errors"
	"strconv"
	"strings"
	"time"

	"shop/config"
	"shop/internal/entity"
	"shop/pkg/logger"
	"shop/pkg/postgres"
)

type CategoryRepo struct {
	pg     *postgres.Postgres
	config *config.Config
	logger *logger.Logger
}

// New -.
func NewCategoryRepo(pg *postgres.Postgres, config *config.Config, logger *logger.Logger) *CategoryRepo {
	return &CategoryRepo{
		pg:     pg,
		config: config,
		logger: logger,
	}
}

func (r *CategoryRepo) Create(ctx context.Context, req *entity.CategoryCreate) error {
	query := `
		INSERT INTO category (
			name
		) VALUES($1)`

	_, err := r.pg.Pool.Exec(ctx, query,
		req.Name,
	)
	if err != nil {
		return err
	}

	return nil
}

func (r *CategoryRepo) GetById(ctx context.Context, req *entity.ById) (*entity.CategoryRes, error) {

	var res entity.CategoryRes
	var createdAt time.Time

	query := `
	SELECT
		id,
		name,
		created_at
	FROM 
		category
	WHERE 
		deleted_at = 0
	AND 
		id = $1
	`

	row := r.pg.Pool.QueryRow(ctx, query, req.Id)
	err := row.Scan(
		&res.Id,
		&res.Name,
		&createdAt,
	)
	if err != nil {
		return nil, err
	}
	res.CreatedAt = createdAt.Format("2006-01-02 15:04:05")

	return &res, nil
}

func (r *CategoryRepo) GetAll(ctx context.Context, req *entity.Filter) (*entity.CategoryGetAllRes, error) {

	resp := &entity.CategoryGetAllRes{}

	query := `
	SELECT
		COUNT(id) OVER () AS total_count,
		id,
		name,
		created_at
	FROM
		category
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
			return nil, errors.New("category list is empty")
		}
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		res := entity.CategoryRes{}
		var count int
		var createdAt time.Time

		err := rows.Scan(
			&count,
			&res.Id,
			&res.Name,
			&createdAt,
		)
		if err != nil {
			return nil, err
		}
		res.CreatedAt = createdAt.Format("2006-01-02 15:04:05")

		resp.Categorys = append(resp.Categorys, res)
		resp.Count = count
	}

	return resp, nil
}

func (r *CategoryRepo) Update(ctx context.Context, req *entity.CategoryUpdate) error {
	query := `
	UPDATE
		category
	SET`

	var conditions []string
	var args []interface{}

	if req.Name != "" && req.Name != "string" {
		conditions = append(conditions, " name = $"+strconv.Itoa(len(args)+1))
		args = append(args, req.Name)
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

func (r *CategoryRepo) Delete(ctx context.Context, req *entity.ById) error {

	_, err := r.pg.Pool.Exec(ctx, `UPDATE category SET deleted_at = EXTRACT(EPOCH FROM NOW()) WHERE id = $1`, req.Id)
	if err != nil {
		return err
	}

	return nil
}

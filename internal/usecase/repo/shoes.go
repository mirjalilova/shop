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

type ShoesRepo struct {
	pg     *postgres.Postgres
	config *config.Config
	logger *logger.Logger
}

// New -.
func NewShoesRepo(pg *postgres.Postgres, config *config.Config, logger *logger.Logger) *ShoesRepo {
	return &ShoesRepo{
		pg:     pg,
		config: config,
		logger: logger,
	}
}

func (r *ShoesRepo) Create(ctx context.Context, req *entity.ShoesCreate) error {
	query := `
		INSERT INTO shoes (
			name,
			size,
			color,
			img_url,
			price,
			description,
			category_id
		) VALUES($1, $2, $3, $4, $5, $6, $7)`

	_, err := r.pg.Pool.Exec(ctx, query,
		req.Name,
		req.Size,
		req.Color,
		req.ImgUrl,
		req.Price,
		req.Description,
		req.CategoryId,
	)
	if err != nil {
		return err
	}

	return nil
}

func (r *ShoesRepo) GetById(ctx context.Context, req *entity.ById) (*entity.ShoesRes, error) {

	var res entity.ShoesRes
	var createdAt time.Time

	query := `
	SELECT
		id,
		name,
		size,
		color,
		img_url,
		price,
		description,
		category_id,
		created_at
	FROM
		shoes
	WHERE
		deleted_at = 0
	AND 
		id = $1
	`

	row := r.pg.Pool.QueryRow(ctx, query, req.Id)
	err := row.Scan(
		&res.Id,
		&res.Name,
		&res.Size,
		&res.Color,
		&res.ImgUrl,
		&res.Price,
		&res.Description,
		&res.CategoryId,
		&createdAt,
	)
	if err != nil {
		return nil, err
	}
	res.CreatedAt = createdAt.Format("2006-01-02 15:04:05")

	return &res, nil
}

func (r *ShoesRepo) GetAll(ctx context.Context, req *entity.Filter) (*entity.ShoesGetAllRes, error) {

	resp := &entity.ShoesGetAllRes{}

	query := `
	SELECT
		COUNT(id) OVER () AS total_count,
		id,
		name,
		size,
		color,
		img_url,
		price,
		description,
		category_id,
		created_at
	FROM
		shoes
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
		res := entity.ShoesRes{}
		var count int
		var createdAt time.Time

		err := rows.Scan(
			&count,
			&res.Id,
			&res.Name,
			&res.Size,
			&res.Color,
			&res.ImgUrl,
			&res.Price,
			&res.Description,
			&res.CategoryId,
			&createdAt,
		)
		if err != nil {
			return nil, err
		}
		res.CreatedAt = createdAt.Format("2006-01-02 15:04:05")

		resp.Shoess = append(resp.Shoess, res)
		resp.Count = count
	}

	return resp, nil
}

func (r *ShoesRepo) Update(ctx context.Context, req *entity.ShoesUpdate) error {
	query := `
	UPDATE
		shoes
	SET`

	var conditions []string
	var args []interface{}

	if req.Name != "" && req.Name != "string" {
		conditions = append(conditions, " name = $"+strconv.Itoa(len(args)+1))
		args = append(args, req.Name)
	}
	if len(req.Size) > 0 {
		conditions = append(conditions, " size = $"+strconv.Itoa(len(args)+1))
		args = append(args, req.Size)
	}
	if len(req.Color) > 0 {
		conditions = append(conditions, " color = $"+strconv.Itoa(len(args)+1))
		args = append(args, req.Color)
	}
	if len(req.ImgUrl) > 0 {
		conditions = append(conditions, " img_url = $"+strconv.Itoa(len(args)+1))
		args = append(args, req.ImgUrl)
	}
	if req.Price > 0 {
		conditions = append(conditions, " price = $"+strconv.Itoa(len(args)+1))
		args = append(args, req.Price)
	}
	if req.Description != "" && req.Description != "string" {
		conditions = append(conditions, " description = $"+strconv.Itoa(len(args)+1))
		args = append(args, req.Description)
	}
	if req.CategoryId != "" && req.CategoryId != "string" {
		conditions = append(conditions, " category_id = $"+strconv.Itoa(len(args)+1))
		args = append(args, req.CategoryId)
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

func (r *ShoesRepo) Delete(ctx context.Context, req *entity.ById) error {

	_, err := r.pg.Pool.Exec(ctx, `UPDATE shoes SET deleted_at = EXTRACT(EPOCH FROM NOW()) WHERE id = $1`, req.Id)
	if err != nil {
		return err
	}

	return nil
}

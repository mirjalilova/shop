package repo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"shop/config"
	"shop/internal/entity"
	"shop/pkg/logger"
	"shop/pkg/postgres"
)

type ProductRepo struct {
	pg     *postgres.Postgres
	config *config.Config
	logger *logger.Logger
}

// New -.
func NewProductRepo(pg *postgres.Postgres, config *config.Config, logger *logger.Logger) *ProductRepo {
	return &ProductRepo{
		pg:     pg,
		config: config,
		logger: logger,
	}
}

func (r *ProductRepo) Create(ctx context.Context, req *entity.ProductCreate) error {
	query := `
		INSERT INTO products (
			name,
			size,
			type,
			img_url,
			price,
			count,
			description,
			category_id
		) VALUES($1, $2, $3, $4, $5, $6, $7, $8)`

	_, err := r.pg.Pool.Exec(ctx, query,
		req.Name,
		req.Size,
		req.Type,
		req.ImgUrl,
		req.Price,
		req.Count,
		req.Description,
		req.CategoryId,
	)
	if err != nil {
		return err
	}

	return nil
}

func (r *ProductRepo) GetById(ctx context.Context, req *entity.ById) (*entity.ProductRes, error) {

	var res entity.ProductRes
	var createdAt time.Time

	query := `
	SELECT
		id,
		name,
		size,
		type,
		img_url,
		price,
		count,
		description,
		category_id,
		created_at
	FROM
		products
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
		&res.Type,
		&res.ImgUrl,
		&res.Price,
		&res.Count,
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

func (r *ProductRepo) GetAll(ctx context.Context, req *entity.ProductGetAllReq) (*entity.ProductGetAllRes, error) {

	resp := &entity.ProductGetAllRes{}

	query := `
	SELECT
		COUNT(id) OVER () AS total_count,
		id,
		name,
		size,
		type,
		img_url,
		price,
		count,
		description,
		category_id,
		created_at
	FROM
		products
	WHERE
		deleted_at = 0
	`

	var args []interface{}

	if req.CategoryId != "" {
		query += fmt.Sprintf(" AND category_id = $%d", len(args)+1)
		args = append(args, req.CategoryId)
	}

	if req.Search != "" {
		query += fmt.Sprintf(" AND name ILIKE $%d", len(args)+1)
		args = append(args, "%"+req.Search+"%")
	}
	if req.Filter.Limit == 0 {
		query += fmt.Sprintf(" OFFSET $%d", len(args)+1)
		args = append(args, req.Filter.Offset)
	} else {
		query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", len(args)+1, len(args)+2)
		args = append(args, req.Filter.Limit)
		args = append(args, req.Filter.Offset)
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
		res := entity.ProductRes{}
		var count int
		var createdAt time.Time

		err := rows.Scan(
			&count,
			&res.Id,
			&res.Name,
			&res.Size,
			&res.Type,
			&res.ImgUrl,
			&res.Price,
			&res.Count,
			&res.Description,
			&res.CategoryId,
			&createdAt,
		)
		if err != nil {
			return nil, err
		}
		res.CreatedAt = createdAt.Format("2006-01-02 15:04:05")

		resp.Products = append(resp.Products, res)
		resp.Count = count
	}

	return resp, nil
}

func (r *ProductRepo) Update(ctx context.Context, req *entity.ProductUpdate) error {
	query := `
	UPDATE
		products
	SET`

	var conditions []string
	var args []interface{}

	if req.Name != "" && req.Name != "string" {
		conditions = append(conditions, " name = $"+strconv.Itoa(len(args)+1))
		args = append(args, req.Name)
	}
	if req.Size > 0 {
		conditions = append(conditions, " size = $"+strconv.Itoa(len(args)+1))
		args = append(args, req.Size)
	}
	if req.Type > "" && req.Type != "string" {
		conditions = append(conditions, " color = $"+strconv.Itoa(len(args)+1))
		args = append(args, req.Type)
	}
	if req.ImgUrl > "" && req.ImgUrl != "string" {
		conditions = append(conditions, " img_url = $"+strconv.Itoa(len(args)+1))
		args = append(args, req.ImgUrl)
	}
	if req.Price > 0 {
		conditions = append(conditions, " price = $"+strconv.Itoa(len(args)+1))
		args = append(args, req.Price)
	}
	if req.Count > 0 {
		conditions = append(conditions, " count = $"+strconv.Itoa(len(args)+1))
		args = append(args, req.Count)
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

func (r *ProductRepo) Delete(ctx context.Context, req *entity.ById) error {

	_, err := r.pg.Pool.Exec(ctx, `UPDATE products SET deleted_at = EXTRACT(EPOCH FROM NOW()) WHERE id = $1`, req.Id)
	if err != nil {
		return err
	}

	return nil
}

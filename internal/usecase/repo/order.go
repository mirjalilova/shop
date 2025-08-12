package repo

import (
	"context"
	"fmt"
	"shop/config"
	"shop/internal/entity"
	"shop/pkg/logger"
	"shop/pkg/postgres"
)

type OrderRepo struct {
	pg     *postgres.Postgres
	config *config.Config
	logger *logger.Logger
}

func NewOrderRepo(pg *postgres.Postgres, config *config.Config, logger *logger.Logger) *OrderRepo {
	return &OrderRepo{
		pg:     pg,
		config: config,
		logger: logger,
	}
}

func (r *OrderRepo) Create(ctx context.Context, req *entity.OrderCreate) error {

	tr, err := r.pg.Pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("error while begin transactions")
	}

	loc := fmt.Sprintf("POINT(%f %f)", req.Location.Longitude, req.Location.Latitude)

	query := `
		INSERT INTO orders (
			user_id,
			bucket_id,
			status,
			location,
			description,
			payment_type
		) VALUES($1, $2, $3, $4, $5, $6)
		`

	_, err = tr.Exec(ctx, query,
		req.UserID,
		req.BucketID,
		req.Status,
		loc,
		req.Description,
		req.PaymentType,
	)
	if err != nil {
		tr.Rollback(ctx)
		return err
	}

	query = `
		SELECT
			bi.product_id,
			bi.count
		FROM
			buckets b 
		JOIN 
			bucket_item bi ON b.id = bi.bucket_id
		WHERE
			b.id = $1
		AND 
			b.deleted_at = 0
		`
	rows, err := r.pg.Pool.Query(ctx, query, req.BucketID)
	if err != nil {
		tr.Rollback(ctx)
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var count int
		var product_id string
		err = rows.Scan(
			&product_id,
			&count,
		)
		if err != nil {
			tr.Rollback(ctx)
			return err
		}

		query := `
			UPDATE products SET count = count - $2 WHERE id = $1`
		_, err = tr.Exec(ctx, query, product_id, count)
		if err != nil {
			tr.Rollback(ctx)
			return err
		}
	}

	query = `
		UPDATE buckets SET status = true WHERE id = $2`
	_, err = tr.Exec(ctx, query, req.BucketID)
	if err != nil {
		tr.Rollback(ctx)
		return err
	}

	err = tr.Commit(ctx)
	if err != nil {
		tr.Commit(ctx)
		return err
	}
	return nil
}

func (r *OrderRepo) GetOrders(ctx context.Context, status string, user_id string) (*[]entity.OrderRes, error) {

	res := []entity.OrderRes{}

	query := `
		SELECT
			id,
			status,
			ST_Y(location) AS latitude,
			ST_X(location) AS longitude,
			description,
			payment_type,
			bucket_id
		FROM
			orders
		WHERE
			deleted_at = 0
		`

	if status != "" {
		query += fmt.Sprintf(" AND status = '%s'", status)
	}

	if user_id != "" {
		query += fmt.Sprintf(" AND user_id = '%s'", user_id)
	}

	rows, err := r.pg.Pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		order := &entity.OrderRes{}
		var bucket_id string
		err = rows.Scan(
			&order.ID,
			&order.Status,
			&order.Location.Latitude,
			&order.Location.Longitude,
			&order.Description,
			&order.PaymentType,
			&bucket_id,
		)
		if err != nil {
			return nil, err
		}

		query = `
		SELECT
			bi.id,
			bi.product_id,
			p.name,
			p.size,
			p.type,
			p.price as product_price,
			p.img_url,
			bi.count,
			bi.price,
			b.total_price
		FROM
			buckets b 
		JOIN bucket_item bi ON b.id = bi.bucket_id
		JOIN products p ON bi.product_id = p.id
		WHERE
			b.id = $1
		AND 
			b.deleted_at = 0
		`
		rows, err = r.pg.Pool.Query(ctx, query, bucket_id)
		if err != nil {
			return nil, err
		}
		defer rows.Close()

		items := []entity.OrderItemRes{}
		for rows.Next() {
			item := &entity.OrderItemRes{}
			err = rows.Scan(
				&item.Id,
				&item.ProductID,
				&item.ProductName,
				&item.ProductSize,
				&item.ProductType,
				&item.ProductPrice,
				&item.ImgUrl,
				&item.Count,
				&item.Price,
				&order.TotalPrice,
			)
			if err != nil {
				return nil, err
			}
			items = append(items, *item)
		}

		order.Orders = items

		res = append(res, *order)
	}

	return &res, nil
}

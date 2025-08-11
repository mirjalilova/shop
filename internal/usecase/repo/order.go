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
		req.Location,
		req.Description,
		req.PaymentType,
	)
	if err != nil {
		tr.Rollback(ctx)
		return err
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

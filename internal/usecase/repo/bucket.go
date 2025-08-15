package repo

import (
	"context"
	"fmt"
	"shop/config"
	"shop/internal/entity"
	"shop/pkg/logger"
	"shop/pkg/postgres"
	"strings"
)

type BucketRepo struct {
	pg     *postgres.Postgres
	config *config.Config
	logger *logger.Logger
}

func NewBucketRepo(pg *postgres.Postgres, config *config.Config, logger *logger.Logger) *BucketRepo {
	return &BucketRepo{
		pg:     pg,
		config: config,
		logger: logger,
	}
}

func (r *BucketRepo) Create(ctx context.Context, req *entity.BucketItemCreate) error {

	var bucket_id string

	tr, err := r.pg.Pool.Begin(ctx)
	if err != nil {
		return err
	}

	query := "SELECT id from buckets where user_id = $1 AND deleted_at = 0 AND status = false LIMIT 1"

	err = tr.QueryRow(ctx, query, req.UserID).Scan(&bucket_id)
	if err != nil {
		tr.Rollback(ctx)
		return err
	}

	var price float64
	query = `
		INSERT INTO bucket_item (bucket_id, product_id, count, price)
		SELECT
			$1 AS bucket_id,
			$2 AS product_id,
			$3 AS count,
			price
		FROM products
		WHERE id = $2
		RETURNING price`

	err = tr.QueryRow(ctx, query,
		bucket_id,
		req.ProductID,
		req.Count,
	).Scan(&price)
	if err != nil {
		tr.Rollback(ctx)
		return err
	}

	query = `
		UPDATE buckets  
		SET total_price = total_price + $1 WHERE id = $2`

	_, err = tr.Exec(ctx, query, float64(req.Count)*price, bucket_id)
	if err != nil {
		tr.Rollback(ctx)
		return err
	}

	return tr.Commit(ctx)
}

func (r *BucketRepo) GetBucket(ctx context.Context, user_id string) (*entity.BucketRes, error) {
	res := &entity.BucketRes{}

	query := `
		SELECT
			b.id,
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
			b.status = false 
		AND
			b.user_id = $1
		AND 
			b.deleted_at = 0
		`
	rows, err := r.pg.Pool.Query(ctx, query, user_id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := []entity.BucketItemRes{}
	for rows.Next() {
		item := &entity.BucketItemRes{}
		var total_price float32
		var bucket_id string
		err = rows.Scan(
			&bucket_id,
			&item.Id,
			&item.ProductID,
			&item.ProductName,
			&item.ProductSize,
			&item.ProductType,
			&item.ProductPrice,
			&item.ImgUrl,
			&item.Count,
			&item.Price,
			&total_price,
		)
		if err != nil {
			return nil, err
		}
		items = append(items, *item)
		res.TotalPrice = total_price
		res.BucketID = bucket_id
	}

	res.Buskets = items

	return res, nil
}

func (r *BucketRepo) Update(ctx context.Context, req *entity.BucketItemUpdate) error {
	tr, err := r.pg.Pool.Begin(ctx)
	if err != nil {
		return err
	}

	var (
		conditions []string
		args       []interface{}
	)

	if req.Count != 0 {
		conditions = append(conditions, fmt.Sprintf("count = $%d", len(args)+1))
		args = append(args, req.Count)
	}

	if req.Price != 0 {
		conditions = append(conditions, fmt.Sprintf("price = $%d", len(args)+1))
		args = append(args, req.Price)
	}

	conditions = append(conditions, "updated_at = now()")

	query := fmt.Sprintf(`
		UPDATE bucket_item
		SET %s
		WHERE id = $%d AND deleted_at = 0`,
		strings.Join(conditions, ", "),
		len(args)+1,
	)

	args = append(args, req.Id) 

	if _, err := tr.Exec(ctx, query, args...); err != nil {
		tr.Rollback(ctx)
		return err
	}

	var bucketID string
	if err := tr.QueryRow(ctx,
		`SELECT bucket_id FROM bucket_item WHERE id = $1`,
		req.Id,
	).Scan(&bucketID); err != nil {
		tr.Rollback(ctx)
		return err
	}

	rows, err := tr.Query(ctx, `
		SELECT count, price
		FROM bucket_item
		WHERE bucket_id = $1 AND deleted_at = 0`,
		bucketID,
	)
	if err != nil {
		tr.Rollback(ctx)
		return err
	}
	defer rows.Close()

	var totalPrice float64
	for rows.Next() {
		var count int
		var price float64
		if err := rows.Scan(&count, &price); err != nil {
			tr.Rollback(ctx)
			return err
		}
		totalPrice += float64(count) * price
	}

	if _, err := tr.Exec(ctx,
		`UPDATE buckets SET total_price = $1 WHERE id = $2`,
		totalPrice, bucketID,
	); err != nil {
		tr.Rollback(ctx)
		return err
	}

	return tr.Commit(ctx)
}

func (r *BucketRepo) Delete(ctx context.Context, id string) error {
	query := `
		DELETE FROM 
			bucket_item
		WHERE id = $1 AND deleted_at = 0`

	_, err := r.pg.Pool.Exec(ctx, query, id)
	if err != nil {
		return err
	}

	return nil
}

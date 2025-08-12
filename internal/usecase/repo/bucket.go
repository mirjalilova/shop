package repo

import (
	"context"
	"shop/config"
	"shop/internal/entity"
	"shop/pkg/logger"
	"shop/pkg/postgres"
	"strconv"
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

	query := "SELECT id from buckets where user_id = $1 AND deleted_at = 0"

	err := r.pg.Pool.QueryRow(ctx, query, req.UserID).Scan(&bucket_id)
	if err != nil {
		return err
	}

	query = `
		INSERT INTO bucket_item (bucket_id, product_id, count, price)
		SELECT
			$1 AS bucket_id,
			$2 AS product_id,
			$3 AS count,
			price
		FROM products
		WHERE id = $2`

	_, err = r.pg.Pool.Exec(ctx, query,
		bucket_id,
		req.ProductID,
		req.Count,
	)
	if err != nil {
		return err
	}

	return nil
}

func (r *BucketRepo) GetBucket(ctx context.Context, user_id string) (*entity.BucketRes, error) {
	res := &entity.BucketRes{}

	query := `
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
			&total_price,
		)
		if err != nil {
			return nil, err
		}
		items = append(items, *item)
		res.TotalPrice = total_price
	}

	res.Buskets = items

	return res, nil
}

func (r *BucketRepo) Update(ctx context.Context, req *entity.BucketItemUpdate) error {
	query := `
		UPDATE 
			bucket_item
		SET `

	var conditions []string
	var args []interface{}

	if req.Count != 0 {
		conditions = append(conditions, " count = $"+strconv.Itoa(len(args)+1))
		args = append(args, req.Count)
	}

	if req.Price != 0 {
		conditions = append(conditions, " price = $"+strconv.Itoa(len(args)+1))
		args = append(args, req.Price)
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

package repo

import (
	"context"
	"fmt"
	"shop/config"
	"shop/internal/entity"
	"shop/pkg/logger"
	"shop/pkg/postgres"
	"shop/pkg/telegram"
	"time"
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
			location,
			description,
			payment_type
		) VALUES($1, $2, $3, $4, $5)
		RETURNING id
		`

	var orderID string
	err = tr.QueryRow(ctx, query,
		req.UserID,
		req.BucketID,
		loc,
		req.Description,
		req.PaymentType,
	).Scan(&orderID)
	if err != nil {
		tr.Rollback(ctx)
		return fmt.Errorf("error while creating order: %w", err)
	}

	query = `
		SELECT bi.product_id, bi.count
		FROM buckets b 
		JOIN bucket_item bi ON b.id = bi.bucket_id
		WHERE b.id = $1 AND b.deleted_at = 0
	`
	rows, err := r.pg.Pool.Query(ctx, query, req.BucketID)
	if err != nil {
		tr.Rollback(ctx)
		return fmt.Errorf("error while querying bucket items: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var count int
		var productID string
		if err := rows.Scan(&productID, &count); err != nil {
			tr.Rollback(ctx)
			return err
		}

		_, err = tr.Exec(ctx, `UPDATE products SET count = count - $2 WHERE id = $1`, productID, count)
		if err != nil {
			tr.Rollback(ctx)
			return fmt.Errorf("error while updating product count: %w", err)
		}
	}

	_, err = tr.Exec(ctx, `UPDATE buckets SET status = true WHERE id = $1`, req.BucketID)
	if err != nil {
		tr.Rollback(ctx)
		return fmt.Errorf("error while updating bucket status: %w", err)
	}

	_, err = tr.Exec(ctx, `INSERT INTO buckets (user_id)VALUES($1)`, req.UserID)
	if err != nil {
		tr.Rollback(ctx)
		return fmt.Errorf("error while creating bucket: %w", err)
	}

	if err := tr.Commit(ctx); err != nil {
		tr.Rollback(ctx)
		return fmt.Errorf("error while committing transaction: %w", err)
	}

	go r.sendOrderToTelegram(context.Background(), orderID)

	return nil
}

func (r *OrderRepo) GetOrders(ctx context.Context, status string, user_id string) (*[]entity.OrderRes, error) {

	res := []entity.OrderRes{}

	query := `
		SELECT
			id,
			status,
			ST_Y(location::geometry) AS latitude,
			ST_X(location::geometry) AS longitude,
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

	orderRows, err := r.pg.Pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer orderRows.Close()

	for orderRows.Next() {
		order := entity.OrderRes{}
		var bucket_id string
		err = orderRows.Scan(
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

		itemQuery := `
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
		itemRows, err := r.pg.Pool.Query(ctx, itemQuery, bucket_id)
		if err != nil {
			return nil, err
		}

		items := []entity.OrderItemRes{}
		for itemRows.Next() {
			var item entity.OrderItemRes
			err = itemRows.Scan(
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
				itemRows.Close()
				return nil, err
			}
			items = append(items, item)
		}
		itemRows.Close()

		order.Orders = items
		res = append(res, order)
	}

	return &res, nil
}

func (r *OrderRepo) Update(ctx context.Context, status, id string) error {
	query := `
		UPDATE orders
		SET status = $1
		WHERE id = $2 AND deleted_at = 0
	`
	_, err := r.pg.Pool.Exec(ctx, query, status, id)
	return err
}

func (r *OrderRepo) sendOrderToTelegram(ctx context.Context, orderID string) {
	query := `
		SELECT 
			o.id, 
			u.name, 
			u.phone_number, 
			o.description, 
			o.payment_type, 
			o.created_at,
		    ST_Y(o.location::geometry) AS latitude,
			ST_X(o.location::geometry) AS longitude,
			b.total_price
		FROM orders o
		JOIN users u ON o.user_id = u.id
		JOIN buckets b ON o.bucket_id = b.id
		WHERE o.id = $1
	`
	row := r.pg.Pool.QueryRow(ctx, query, orderID)

	var (
		id, name, phone, description, paymentType string
		createdAt                                 time.Time
		totalPrice                                int
	)
	lat := entity.Location{}
	if err := row.Scan(&id, &name, &phone, &description, &paymentType, &createdAt, &lat.Latitude, &lat.Longitude, &totalPrice); err != nil {
		r.logger.Error("failed to fetch order for telegram: ", err)
		return
	}

	itemRows, err := r.pg.Pool.Query(ctx, `
		SELECT 
			p.name, 
			bi.count, 
			bi.price
		FROM 
			bucket_item bi
		JOIN 
			products p ON bi.product_id = p.id
		WHERE 
			bi.bucket_id = (SELECT bucket_id FROM orders WHERE id = $1)
	`, orderID)
	if err != nil {
		r.logger.Error("failed to fetch order items: ", err)
		return
	}
	defer itemRows.Close()

	itemsText := ""
	for itemRows.Next() {
		var pname string
		var count int
		var price float64
		if err := itemRows.Scan(&pname, &count, &price); err != nil {
			r.logger.Error("scan item error: ", err)
			return
		}
		itemsText += fmt.Sprintf("• %s x%d — %.2f\n", pname, count, price)
	}

	message := fmt.Sprintf(
		"<b>🆕 Yangi Buyurtma</b>\n\n🆔 ID: %s\n👤 Mijoz: %s\n📞 Telefon: %s\n📍 Joylashuv: %.6f, %.6f\n💳 To‘lov turi: %s\n🛒 Buyurtmalar:\n%s\n💰 Jami: %d\n🕒 Sana: %s",
		id, name, phone, lat.Latitude, lat.Longitude, paymentType, itemsText, totalPrice, createdAt.Format("2006-01-02 15:04:05"),
	)

	tg := telegram.NewClient(r.config.Telegram.Token, r.config.Telegram.ChatID)
	if err := tg.SendMessage(message); err != nil {
		r.logger.Error("telegram send error: ", err)
	}
}
